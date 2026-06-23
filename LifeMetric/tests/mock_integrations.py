#!/usr/bin/env python3
"""
LifeMetrics Global Integration Mock Server (Zero-Dependency)
Author: Hana (Backend Engineer)

This server simulates our external healthcare and billing integrations:
1. HL7 FHIR R4 EHR Server (Epic, Cerner, PointClickCare) to sync patient contexts & contactless vitals.
2. mybalanceapp.com Subscription Portal to simulate SaaS subscription upgrades and payment verify hooks.

Runs on local port 8081.
"""

import json
from http.server import HTTPServer, BaseHTTPRequestHandler

PORT = 8081

class MockIntegrationHandler(BaseHTTPRequestHandler):
    
    def log_message(self, format, *args):
        # Override to suppress standard HTTP logging to keep stdout clean
        return

    def _set_headers(self, status_code=200, content_type="application/json"):
        self.send_response(status_code)
        self.send_header("Content-Type", content_type)
        self.send_header("Access-Control-Allow-Origin", "*")
        self.send_header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        self.send_header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        self.end_headers()

    def do_OPTIONS(self):
        # Handle CORS preflight requests
        self._set_headers(204)

    def do_GET(self):
        # 1. Handle FHIR Patient Query (Epic/Cerner sandbox context)
        # e.g., GET /fhir/r4/Patient/102948
        if self.path == "/fhir/r4/Patient/102948":
            patient_payload = {
                "resourceType": "Patient",
                "id": "102948",
                "active": True,
                "name": [{
                    "use": "official",
                    "family": "Jenkins",
                    "given": ["Sarah"]
                }],
                "gender": "female",
                "birthDate": "1948-03-12",
                "extension": [
                    {
                        "url": "http://lifemetrics.com/fhir/StructureDefinition/mobility-index",
                        "valueString": "Uses Cane / Shuffling Gait Baseline"
                    },
                    {
                        "url": "http://lifemetrics.com/fhir/StructureDefinition/patient-weight-kg",
                        "valueInteger": 64
                    }
                ]
            }
            self._set_headers(200)
            self.wfile.write(json.dumps(patient_payload, indent=2).encode('utf-8'))
            print("\n[MOCK EHR FHIR] 🔎 Received Patient Context Query for Patient #102948 (Sarah Jenkins)")
            print("  ↳ Returned: Female, Born 1948-03-12, Uses Cane / Shuffling Gait (64kg)")
            return

        # 2. Handle mybalanceapp.com plan verification status
        # e.g., GET /mybalanceapp/subscriptions/verify?tenant_id=...
        if "/mybalanceapp/subscriptions/verify" in self.path:
            verify_payload = {
                "tenant_id": "aa1c305b-9d4f-4d67-8e12-3cf9052735ba",
                "subscription_status": "premium_active",
                "billing_cycle_ends": "2026-07-23T00:00:00Z"
            }
            self._set_headers(200)
            self.wfile.write(json.dumps(verify_payload, indent=2).encode('utf-8'))
            print("\n[MOCK BILLING] 💳 Verified subscription status for Tenant Facility")
            print("  ↳ Returned: premium_active, Cycle Ends 2026-07-23")
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

        # 1. Handle FHIR Observation Ingestion (Vitals Push)
        if self.path == "/fhir/r4/Observation":
            # Extract heart rate and respiratory rate values from FHIR payload components
            heart_rate = None
            resp_rate = None
            
            if body.get("resourceType") == "Observation":
                components = body.get("component", [])
                for comp in components:
                    loinc_code = comp.get("code", {}).get("coding", [{}])[0].get("code")
                    val = comp.get("valueQuantity", {}).get("value")
                    if loinc_code == "8867-4":
                        heart_rate = val
                    elif loinc_code == "9279-1":
                        resp_rate = val

            print("\n[MOCK EHR FHIR] 📈 Received Contactless Vital Signs Observation Write:")
            print(f"  ↳ Patient ID: {body.get('subject', {}).get('reference')}")
            print(f"  ↳ Ingestion Timestamp: {body.get('effectiveDateTime')}")
            if heart_rate and resp_rate:
                print(f"  ↳ Extracted Biometrics: Heart Rate = {heart_rate} BPM | Respiratory Rate = {resp_rate} Br/min")
            else:
                print("  ↳ Raw Data:", json.dumps(body, indent=2))
                
            self._set_headers(201)
            self.wfile.write(json.dumps({"status": "Observation Created", "fhir_id": "obs_9028475a"}, indent=2).encode('utf-8'))
            return

        # 2. Handle mybalanceapp.com Checkout
        if self.path == "/mybalanceapp/checkout":
            print("\n[MOCK BILLING] 🛍️ Processing Subscription Checkout via mybalanceapp.com API:")
            print(f"  ↳ Tenant Business: {body.get('business_name', 'Unnamed Facility')}")
            print(f"  ↳ Selected Plan: {body.get('selected_plan', 'premium')}")
            print(f"  ↳ Payment Method: {body.get('payment_method', 'card')}")
            print("  ↳ Processing Secure 3D-Secure Verification Gate...")
            
            checkout_response = {
                "transaction_id": "tx_mybal_928475a0",
                "status": "success",
                "charge_amount_gbp": 20.00,
                "timestamp": "2026-06-23T15:20:00Z"
            }
            self._set_headers(200)
            self.wfile.write(json.dumps(checkout_response, indent=2).encode('utf-8'))
            return

        self._set_headers(404)
        self.wfile.write(json.dumps({"error": "Endpoint not found"}).encode('utf-8'))

def run(server_class=HTTPServer, handler_class=MockIntegrationHandler):
    server_address = ('', PORT)
    httpd = server_class(server_address, handler_class)
    print("=" * 70)
    print(f"   LIFEMETRICS EHR (FHIR) & mybalanceapp.com MOCK INTEGRATIONS SERVER")
    print(f"   Running on local port: {PORT}")
    print("=" * 70)
    print("Press Ctrl+C to shut down the server.\n")
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        print("\nIntegration mock server shut down.")
        httpd.server_close()

if __name__ == '__main__':
    run()
