#!/usr/bin/env python3
"""
LifeMetrics Ingestion API Gateway Mock Server (Zero-Dependency)
Author: Hana (Backend Engineer)

This server acts as a drop-in replacement for our Go API Gateway on port 8080.
It allows full-stack testing (Simulator -> Mock Gateway -> Clinician Dashboard)
for users who do not have the Go development compiler installed locally.

Supports CORS and exposes:
- POST /api/v1/devices/register (Device pairing)
- POST /api/v1/devices/calibrate (Calibration trigger)
- POST /api/v1/events/ingest (Low-latency stream ingestion & alert evaluation)
- GET /api/v1/alerts/poll (Local polling fallback for the browser dashboard)
"""

import json
import uuid
from http.server import HTTPServer, BaseHTTPRequestHandler

PORT = 8080

# Simple global in-memory queue to hold active alerts
triggered_alerts_queue = []
registered_devices = {}

class MockGatewayHandler(BaseHTTPRequestHandler):

    def log_message(self, format, *args):
        # Override to suppress standard HTTP logging to keep console clean
        return

    def _set_headers(self, status_code=200, content_type="application/json"):
        self.send_response(status_code)
        self.send_header("Content-Type", content_type)
        # Enable CORS for local testing
        self.send_header("Access-Control-Allow-Origin", "*")
        self.send_header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        self.send_header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        self.end_headers()

    def do_OPTIONS(self):
        # Handle CORS preflight requests
        self._set_headers(204)

    def do_GET(self):
        # Expose health route
        if self.path == "/health" or self.path == "/":
            self._set_headers(200)
            self.wfile.write(json.dumps({"status": "healthy", "gateway": "Python Mock Gateway"}, indent=2).encode('utf-8'))
            return

        # Expose alert polling route for the dashboard
        if self.path == "/api/v1/alerts/poll":
            self._set_headers(200)
            if len(triggered_alerts_queue) > 0:
                # Pop the oldest alert from the queue
                alert = triggered_alerts_queue.pop(0)
                self.wfile.write(json.dumps(alert).encode('utf-8'))
            else:
                self.wfile.write(json.dumps({}).encode('utf-8'))
            return

        self._set_headers(404)
        self.wfile.write(json.dumps({"error": "Endpoint not found"}).encode('utf-8'))

    def do_POST(self):
        content_length = int(self.headers['Content-Length'])
        post_data = self.rfile.read(content_length).decode('utf-8')
        
        try:
            body = json.loads(post_data) if post_data else {}
        except json.JSONDecodeError:
            self._set_headers(400)
            self.wfile.write(json.dumps({"error": "Malformed JSON payload"}).encode('utf-8'))
            return

        # 1. Device Registration Endpoint
        if self.path == "/api/v1/devices/register":
            device_uuid = body.get("device_uuid") or str(uuid.uuid4())
            room_label = body.get("room_label", "Default Room")
            
            registered_devices[device_uuid] = {
                "device_uuid": device_uuid,
                "room_label": room_label,
                "status": "registered"
            }
            
            print(f"\n[GATEWAY] ➕ Registered sensor device {device_uuid} to {room_label}")
            
            self._set_headers(201)
            self.wfile.write(json.dumps({
                "status": "registered",
                "device_uuid": device_uuid,
                "room_label": room_label,
                "tenant_id": "aa1c305b-9d4f-4d67-8e12-3cf9052735ba"
            }, indent=2).encode('utf-8'))
            return

        # 2. Device Calibration Endpoint
        if self.path == "/api/v1/devices/calibrate":
            device_uuid = body.get("device_uuid")
            if device_uuid in registered_devices:
                registered_devices[device_uuid]["status"] = "calibrating"
                
            print(f"[GATEWAY] ⏳ Initiated 48-hour local room calibration phase for sensor {device_uuid}")
            
            self._set_headers(200)
            self.wfile.write(json.dumps({
                "status": "calibrating",
                "device_uuid": device_uuid,
                "phase_duration": "48 hours",
                "message": "Edge-AI calibration triggered. Mapping static background clutter."
            }, indent=2).encode('utf-8'))
            return

        # 3. Telemetry Ingestion Endpoint
        if self.path == "/api/v1/events/ingest":
            device_uuid = body.get("device_uuid")
            tele_data = body.get("data", {})
            
            v_z = tele_data.get("velocity_z", 0.0)
            is_dangling = tele_data.get("is_dangling", False)
            hr = tele_data.get("vitals_bpm", 74)
            rr = tele_data.get("vitals_breaths", 16)
            
            print(f"📡 Ingested telemetry from {device_uuid} -> HR: {hr} BPM, RR: {rr} Br/m, Position Z: {tele_data.get('com_z')}m")
            
            alert = None
            if v_z < -1.8:
                # Trigger a fall event
                alert = {
                    "event_id": str(uuid.uuid4()),
                    "event_type": "alert.fall.predictive",
                    "device_uuid": device_uuid,
                    "room_label": "Room 204 - Bathroom"
                }
            elif is_dangling:
                # Trigger a bed dangling event
                alert = {
                    "event_id": str(uuid.uuid4()),
                    "event_type": "alert.bed_exit.dangling",
                    "device_uuid": device_uuid,
                    "room_label": "Room 102 - Bed"
                }

            if alert:
                print(f"🚨 [GATEWAY ALERT ENGINE] Triggered {alert['event_type']} for {alert['room_label']}")
                triggered_alerts_queue.append(alert)
                
                self._set_headers(202)
                self.wfile.write(json.dumps({
                    "status": "accepted",
                    "message": "Telemetry received, critical alert dispatched",
                    "triggered_event_id": alert["event_id"],
                    "event_type": alert["event_type"]
                }).encode('utf-8'))
                return

            self._set_headers(202)
            self.wfile.write(json.dumps({
                "status": "accepted",
                "message": "Telemetry metadata ingested successfully"
            }).encode('utf-8'))
            return

        self._set_headers(404)
        self.wfile.write(json.dumps({"error": "Endpoint not found"}).encode('utf-8'))

def run(server_class=HTTPServer, handler_class=MockGatewayHandler):
    server_address = ('', PORT)
    httpd = server_class(server_address, handler_class)
    print("=" * 70)
    print(f"   LIFEMETRICS INGESTION API GATEWAY - MOCK SERVER")
    print(f"   Running on local port: {PORT}")
    print("=" * 70)
    print("Press Ctrl+C to shut down the server.\n")
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        print("\nMock gateway server shut down.")
        httpd.server_close()

if __name__ == '__main__':
    run()
