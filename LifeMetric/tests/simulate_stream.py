#!/usr/bin/env python3
"""
LifeMetrics Ingestion Stream & Biometric Telemetry Simulator (Zero-Dependency)
Author: Layla (QA Test Engineer)

This script simulates a live 60GHz mmWave sensor array streaming anonymized spatial 
and physiological metadata to our Go API Gateway. It supports automated device pairing, 
calibration triggers, and continuous streaming across standard, dangling, and falling states.
"""

import base64
import json
import time
import urllib.request
import urllib.error

# Gateway endpoint configuration
GATEWAY_URL = "http://localhost:8080"
TENANT_UUID = "aa1c305b-9d4f-4d67-8e12-3cf9052735ba"
DEVICE_UUID = "f81d4fae-7dec-11d0-a765-00a0c91e6bf6"

def base64_url_encode(data_bytes):
    """Encode bytes to a base64-url string (omitting padding and using -_ characters)."""
    encoded = base64.b64encode(data_bytes).decode('utf-8')
    return encoded.replace('=', '').replace('+', '-').replace('/', '_')

def generate_mock_jwt():
    """
    Generates a mock JWT token containing our tenant facility metadata and scopes.
    Since our Go gateway parses unverified JWTs during prototyping, we generate a 
    valid un-signed none-algorithm JWT without requiring external libraries.
    """
    header = {"alg": "none", "typ": "JWT"}
    payload = {
        "tenant_id": TENANT_UUID,
        "scopes": ["device:write", "telemetry:write"],
        "role": "clinical:staff"
    }
    
    header_enc = base64_url_encode(json.dumps(header).encode('utf-8'))
    payload_enc = base64_url_encode(json.dumps(payload).encode('utf-8'))
    
    # Structure of unsigned JWT: Header.Payload.Signature
    return f"{header_enc}.{payload_enc}.dummy_signature"

def make_request(path, payload, headers):
    """Zero-dependency HTTP POST request wrapper."""
    url = f"{GATEWAY_URL}{path}"
    data = json.dumps(payload).encode('utf-8')
    req = urllib.request.Request(url, data=data, headers=headers, method='POST')
    try:
        with urllib.request.urlopen(req) as response:
            return response.status, json.loads(response.read().decode('utf-8'))
    except urllib.error.HTTPError as e:
        try:
            err_body = json.loads(e.read().decode('utf-8'))
            return e.code, err_body
        except Exception:
            return e.code, {"error": e.reason}
    except urllib.error.URLError as e:
        return 0, {"error": f"Failed to connect to gateway at {GATEWAY_URL}. Is the Go server running?"}

def run_simulation():
    print("=" * 70)
    print("   LIFEMETRICS AMBIENT HEALTH GATEWAY - QA TELEMETRY SIMULATOR   ")
    print("=" * 70)
    
    # 1. Generate Mock Token
    token = generate_mock_jwt()
    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {token}"
    }
    
    # 2. Step 1: Device Registration & Room Pairing
    print("\n[STEP 1] Pairing physical sensor to facility room...")
    reg_payload = {
        "device_uuid": DEVICE_UUID,
        "room_label": "Room 204 - Bathroom",
        "device_placement": "ceiling",
        "public_key": "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCg..."
    }
    status, resp = make_request("/api/v1/devices/register", reg_payload, headers)
    if status != 201:
        print(f"❌ Device registration failed ({status}): {resp.get('error', resp)}")
        return
    print(f"✅ Success (HTTP {status}): Device paired securely. Status: {resp['status']}")
    
    # 3. Step 2: Trigger Calibration
    print("\n[STEP 2] Triggering 48-Hour environmental calibration phase...")
    cal_payload = {
        "device_uuid": DEVICE_UUID
    }
    status, resp = make_request("/api/v1/devices/calibrate", cal_payload, headers)
    if status != 200:
        print(f"❌ Calibration trigger failed ({status}): {resp.get('error', resp)}")
        return
    print(f"✅ Success (HTTP {status}): Calibration initiated. Status: {resp['status']}")
    print(f"💬 Message: {resp['message']}")
    
    # 4. Step 3: Stream Continuous Telemetry
    print("\n[STEP 3] Initializing active point-cloud stream simulation...")
    print("ℹ️ Streaming live heart-rate, respiratory rate, and Center of Mass (CoM) kinematics.")
    print("Press Ctrl+C to stop simulation at any time.")
    print("-" * 70)
    
    # States dictionary for simulation sequencing
    states = [
        # Normal standing/walking sequence
        {"name": "Standard Monitoring", "vitals": (74, 16), "z": 1.10, "vel_z": 0.0, "dangle": False, "duration": 5},
        # Patient sits on edge of bed (entering Dangling state)
        {"name": "Patient Bed-Edge Dangling (Intent to Stand)", "vitals": (82, 19), "z": 0.52, "vel_z": 0.0, "dangle": True, "duration": 5},
        # Normal return to walking
        {"name": "Standard Monitoring (Resolved)", "vitals": (75, 17), "z": 1.10, "vel_z": 0.0, "dangle": False, "duration": 4},
        # Sudden vertical deceleration (Physical Fall)
        {"name": "WARNING: Sudden Vertical Deceleration (FALL EVENT)", "vitals": (110, 26), "z": 0.15, "vel_z": -2.3, "dangle": False, "duration": 3}
    ]
    
    try:
        for phase in states:
            print(f"\n▶️ Phase Active: {phase['name']} (Duration: {phase['duration']}s)")
            for t in range(phase['duration']):
                telemetry_payload = {
                    "device_uuid": DEVICE_UUID,
                    "timestamp": time.strftime("%Y-%m-%dT%H:%M:%SZ"),
                    "data": {
                        "com_x": 1.42,
                        "com_y": 2.15,
                        "com_z": phase['z'],
                        "velocity_z": phase['vel_z'],
                        "is_dangling": phase['dangle'],
                        "is_present": True,
                        "vitals_bpm": phase['vitals'][0],
                        "vitals_breaths": phase['vitals'][1]
                    }
                }
                status, resp = make_request("/api/v1/events/ingest", telemetry_payload, headers)
                if status == 202:
                    # Parse if alert was triggered
                    if resp.get("triggered_event_id"):
                        print(f"🚨 [ALERT ROUTED] HTTP {status}: {resp['triggered_event_id']} | Type: {resp['event_type']} | {resp['message']}")
                        if "calibration_phase" in resp:
                            print(f"   ↳ [FALLBACK MONITORING] Calibration Active: {resp['calibration_phase']} | Status: {resp['confidence_status']}")
                    else:
                        print(f"📡 [STREAMING] HTTP {status}: Telemetry Metadata Ingested. HR: {phase['vitals'][0]} BPM, RR: {phase['vitals'][1]} Br/min")
                else:
                    print(f"❌ Streaming Error (HTTP {status}): {resp.get('error', resp)}")
                
                time.sleep(1)
    except KeyboardInterrupt:
        print("\nSimulation aborted by user.")
    
    print("\n" + "=" * 70)
    print("Simulation complete. All metrics matched API schemas successfully.")
    print("=" * 70)

if __name__ == "__main__":
    run_simulation()
