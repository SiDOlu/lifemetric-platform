# LifeMetric Agent Memory

## ⚠️ Agent Initialization Mandate
**CRITICAL:** To prevent context drift and misalignment, every agent (human or AI) MUST read the `GEMINI.md` file(s) in the current directory and any relevant subdirectories immediately upon initialization. This is the single source of truth for the current project state, objectives, and agent status. Failure to do so can lead to outdated task assumptions and wasted computational cycles.

## 🎯 Current Objectives
- [x] Task 1: Transition from zero-dependency Python Mock Ingestion to the 
production Go/Fiber API Gateway. (Completed)
- [ ] Task 2: Validate the HL7 FHIR and mybalanceapp.com integration 
server hooks against the Go backend.
- [ ] Task 3: Restrict agent operations to comply with UK v1 scope 
permissions (No direct push/merge to main).
- [x] Task 4: Integrate `calibration_phase` and `confidence_status` parameters into Go gateway ingestion payload and alerts (Hana). (Completed)
- [x] Task 5: Implement WebSocket reconnection logic in `app.js` to fix 
"one-way street" fallback issue. (Completed)
- [ ] Task 6: Validate ESP32-S3 hardware simulation against production Go/Fiber gateway. (In Progress)

## 📊 Current Task Status
- **API Gateway (Go/Fiber):** Core telemetry ingestion, alert dispatching, and new metadata parameters (`calibration_phase`, `confidence_status`) are fully implemented and tested.
- **Frontend Dashboard:** WebSocket + HTTP Polling fallback is functional, with an exponential backoff reconnection mechanism implemented to fix the "one-way street" issue.
- **Mock Integration:** FHIR and Stripe mocks are ready for validation.
- **Hardware Integration:** Transitioning Python telemetry simulator from mock environment to production-parity Go gateway.

## 👥 Core Team & Agent Status
- **SiDOlu (Lead Dev):** Built mock servers, WebSocket synchronization, and frontend dashboard prototypes.
- **SiDOlu (Arch/PRD):** Managed PRD scope, compliance protocols, and system architecture.
- **Backend-Agent (Hana):** [Idle] Task 4 completed. Ready for next integration or hardware bridge testing.
- **Data-Agent (Python):** [Idle] Ready to refine the interactive radar point-cloud simulation script.
- **Frontend-Agent (Mateus):** [Idle] Awaiting implementation of WebSocket reconnection logic.
- **QA-Agent (Layla):** [Idle] Ready to validate FHIR hooks and gateway stability.
- **Hardware-Agent (Lukas):** [Active] Critical path: ESP32-S3 hardware simulation testing against the production gateway.
- **Security/Compliance (Amani):** [Idle] Ready to audit UK v1 scope adherence and investigate branch protection gaps.
- **Privacy (Renata):** [Idle] Ready to review data minimization for health telemetry.

## 📈 Detailed Progress & History
- **Completed:** Successfully executed the memory auto-update confirmation test.
- **Completed:** Built vanilla HTML5/ES6 interactive clinician dashboard prototype.
- **Completed:** Implemented zero-dependency Python Mock Ingestion Gateway (Port 8080).
- **Completed:** Implemented mock integration servers for HL7 FHIR and stripe.
- **Completed:** Integrated `calibration_phase` and `confidence_status` into Go gateway and alerts.
- **Completed:** Frontend HTTP polling fallback implementation.
- **Completed:** Implemented WebSocket reconnection logic in `app.js` with exponential backoff.
- **In Progress:** Validating ESP32-S3 simulation data against the production Go/Fiber v2 gateway.
- **Blocked On:** Testing Go Fiber v2 gateway against the ESP32-S3 radar hardware simulation data.


## 📂 Known Workspace Map
- `./` -> Root context. Holds `.agent_memory.md` and repository baseline `configurations.
- `./backend` -> Go Fiber v2 architecture, token verification engines, and 
mock memory databases.
- `./frontend` -> Clinician Web Dashboard, CSS styles, and fallback polling scripts.
- `./simulators` -> Python 3 edge simulation telemetry engines and 
interactive point-cloud mock scripts.
- `./.gemini/agents` -> Custom agent profiles and system rule 
restrictions.
