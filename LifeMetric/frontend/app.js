// LifeMetrics Web Dashboard Interactivity & Web Audio API Engine
// Author: Mateus (Frontend Engineer)

// Clock & Time Indicator Update
function updateClock() {
    const liveTimeEl = document.getElementById('live-time');
    if (liveTimeEl) {
        const now = new Date();
        liveTimeEl.textContent = now.toTimeString().split(' ')[0];
    }
}
setInterval(updateClock, 1000);

// Web Audio API Synthesizer Context
let audioCtx = null;
let alertOscillator = null;
let alertGainNode = null;
let soundTimer = null;

function initAudioContext() {
    if (!audioCtx) {
        audioCtx = new (window.AudioContext || window.webkitAudioContext)();
    }
}

// Low-level audio synthesis to bypass asset file dependencies (pure code-to-frequency)
function startSiren() {
    initAudioContext();
    stopSiren(); // Ensure clean state

    alertGainNode = audioCtx.createGain();
    alertGainNode.connect(audioCtx.destination);
    alertGainNode.gain.setValueAtTime(0.2, audioCtx.currentTime); // Sane volume limit

    alertOscillator = audioCtx.createOscillator();
    alertOscillator.type = 'sawtooth';
    alertOscillator.frequency.setValueAtTime(800, audioCtx.currentTime);
    alertOscillator.connect(alertGainNode);
    alertOscillator.start();

    // Loop a warble siren sound (High-Low frequency sweeps)
    let isHigh = true;
    soundTimer = setInterval(() => {
        if (alertOscillator && audioCtx.state === 'running') {
            const targetFreq = isHigh ? 1100 : 500;
            alertOscillator.frequency.exponentialRampToValueAtTime(targetFreq, audioCtx.currentTime + 0.35);
            isHigh = !isHigh;
        }
    }, 400);
}

function startChime() {
    initAudioContext();
    stopSiren(); // Ensure clean state

    alertGainNode = audioCtx.createGain();
    alertGainNode.connect(audioCtx.destination);
    alertGainNode.gain.setValueAtTime(0.15, audioCtx.currentTime);

    alertOscillator = audioCtx.createOscillator();
    alertOscillator.type = 'sine';
    alertOscillator.frequency.setValueAtTime(880, audioCtx.currentTime); // High soft chime
    alertOscillator.connect(alertGainNode);
    alertOscillator.start();

    // Pulse chime at intervals (like a heart rate warning beep)
    soundTimer = setInterval(() => {
        if (alertGainNode && audioCtx.state === 'running') {
            alertGainNode.gain.setValueAtTime(0.15, audioCtx.currentTime);
            alertGainNode.gain.exponentialRampToValueAtTime(0.001, audioCtx.currentTime + 0.4);
        }
    }, 1000);
}

function stopSiren() {
    if (soundTimer) {
        clearInterval(soundTimer);
        soundTimer = null;
    }
    if (alertOscillator) {
        try {
            alertOscillator.stop();
        } catch (e) {}
        alertOscillator.disconnect();
        alertOscillator = null;
    }
    if (alertGainNode) {
        alertGainNode.disconnect();
        alertGainNode = null;
    }
}

// UI State Selectors
const roomCard = document.getElementById('room-204');
const roomBadge = document.getElementById('room-204-badge');
const roomNum = document.getElementById('room-204-num');
const metricPos = document.getElementById('metric-pos');
const vitalHR = document.getElementById('vital-hr');
const vitalRR = document.getElementById('vital-rr');
const fallModal = document.getElementById('fall-modal');

const connDot = document.getElementById('connection-dot');
const connStatus = document.getElementById('connection-status');

// Reset to standard safe state
function setStandardState() {
    stopSiren();
    
    // Clear alarms and modals
    fallModal.classList.add('hidden');
    
    // Reset Card Styles
    roomCard.className = "bg-slate-800 border-2 border-slate-700 rounded-xl p-5 shadow-lg flex flex-col justify-between h-48 transition-all relative overflow-hidden";
    roomBadge.className = "bg-blue-900/40 text-blue-400 border border-blue-900/50 px-2.5 py-0.5 rounded-full text-[10px] font-semibold flex items-center space-x-1.5";
    roomBadge.innerHTML = `<span class="w-1.5 h-1.5 bg-blue-400 rounded-full animate-ping"></span><span>Monitoring</span>`;
    roomNum.className = "text-2xl font-black text-blue-500";
    
    // Update live metrics
    metricPos.textContent = "Standing";
    vitalHR.textContent = "74 BPM";
    vitalRR.textContent = "16 /min";
}

// Set Yellow Warning state for Dangling (Bed-Exit Intent)
function setDanglingState(vitals) {
    stopSiren();
    startChime();

    fallModal.classList.add('hidden');

    // Apply Yellow warning classes
    roomCard.className = "bg-amber-950/20 border-2 border-amber-500 rounded-xl p-5 shadow-lg flex flex-col justify-between h-48 transition-all relative overflow-hidden animate-pulse-yellow";
    roomBadge.className = "bg-amber-900/80 text-amber-200 border border-amber-500/50 px-2.5 py-0.5 rounded-full text-[10px] font-semibold flex items-center space-x-1.5";
    roomBadge.innerHTML = `<span class="w-1.5 h-1.5 bg-amber-400 rounded-full animate-ping"></span><span>DANGLING</span>`;
    roomNum.className = "text-2xl font-black text-amber-500";
    
    // Update metrics
    metricPos.textContent = "Sitting on Edge";
    if (vitals) {
        vitalHR.textContent = vitals.bpm + " BPM";
        vitalRR.textContent = vitals.breaths + " /min";
    } else {
        vitalHR.textContent = "82 BPM";
        vitalRR.textContent = "19 /min";
    }
}

// Set Red Danger state for Fall
function setFallState(vitals) {
    stopSiren();
    startSiren();

    // Trigger full screen alarm modal
    fallModal.classList.remove('hidden');

    // Apply Red flashing classes
    roomCard.className = "bg-red-950/20 border-2 border-red-500 rounded-xl p-5 shadow-lg flex flex-col justify-between h-48 transition-all relative overflow-hidden animate-pulse-red";
    roomBadge.className = "bg-red-900/80 text-red-200 border border-red-500/50 px-2.5 py-0.5 rounded-full text-[10px] font-semibold flex items-center space-x-1.5";
    roomBadge.innerHTML = `<span class="w-1.5 h-1.5 bg-red-400 rounded-full animate-ping"></span><span>FALL DETECTED</span>`;
    roomNum.className = "text-2xl font-black text-red-500";

    // Update metrics
    metricPos.textContent = "Fallen on Floor";
    if (vitals) {
        vitalHR.textContent = vitals.bpm + " BPM";
        vitalRR.textContent = vitals.breaths + " /min";
    } else {
        vitalHR.textContent = "110 BPM";
        vitalRR.textContent = "26 /min";
    }
}

// Acknowledge alert, silencing sound and setting Room to "Assistance Dispatched" status
function setAcknowledgedState() {
    stopSiren();
    fallModal.classList.add('hidden');

    // Set Room card to yellow warning indicating we have dispatched help but not fully resolved
    roomCard.className = "bg-amber-950/10 border-2 border-amber-500/60 rounded-xl p-5 shadow-lg flex flex-col justify-between h-48 transition-all relative overflow-hidden";
    roomBadge.className = "bg-slate-900 text-amber-400 border border-amber-500/30 px-2.5 py-0.5 rounded-full text-[10px] font-semibold flex items-center space-x-1.5";
    roomBadge.innerHTML = `<span>Assistance Dispatched</span>`;
    roomNum.className = "text-2xl font-black text-amber-500/80";
    
    metricPos.textContent = "Fallen (Dispatched)";
}

// Button Click Event Listeners
document.getElementById('btn-standard').addEventListener('click', () => setStandardState());
document.getElementById('btn-dangling').addEventListener('click', () => setDanglingState());
document.getElementById('btn-fall').addEventListener('click', () => setFallState());
document.getElementById('btn-ack').addEventListener('click', setAcknowledgedState);

// ============================================================================
// REAL-TIME WEBSOCKET & HTTP POLLING INTEGRATION
// ============================================================================
let ws = null;
let pollingTimer = null;
let reconnectAttempts = 0;
const MAX_RECONNECT_DELAY = 30000;
let reconnectTimeout = null;

function connectWebSocket() {
    // Clear any pending reconnection attempts
    if (reconnectTimeout) {
        clearTimeout(reconnectTimeout);
        reconnectTimeout = null;
    }

    console.log("[WEBSOCKET] Connecting to LifeMetrics Ingestion API Gateway...");
    ws = new WebSocket("ws://localhost:8080/ws/clinical/alerts");

    ws.onopen = () => {
        console.log("[WEBSOCKET] Connected to Live Gateway.");
        connDot.className = "w-2.5 h-2.5 bg-green-500 rounded-full animate-pulse";
        connStatus.textContent = "Connected to Live Gateway";
        
        // SUCCESS: Stop polling and reset reconnection logic
        reconnectAttempts = 0;
        if (pollingTimer) {
            console.log("[WEBSOCKET] Reconnection successful. Stopping HTTP polling.");
            clearInterval(pollingTimer);
            pollingTimer = null;
        }
    };

    ws.onclose = () => {
        console.log("[WEBSOCKET] Connection closed.");
        ws = null;
        
        // If we are not already polling, start polling
        if (!pollingTimer) {
            startHTTPPolling();
        }
        
        scheduleReconnection();
    };

    ws.onerror = (err) => {
        console.warn("[WEBSOCKET] Connection error occurred.");
        // ws.close() will trigger onclose
        if (ws) ws.close();
    };

    ws.onmessage = (event) => {
        try {
            const alert = JSON.parse(event.data);
            console.log("[WEBSOCKET] Real-Time Alert Received:", alert);
            handleIncomingAlert(alert);
        } catch (e) {
            console.error("[WEBSOCKET] Failed to parse alert data packet:", e);
        }
    };
}

function scheduleReconnection() {
    if (reconnectTimeout) return;

    const delay = Math.min(1000 * Math.pow(2, reconnectAttempts), MAX_RECONNECT_DELAY);
    console.log(`[WEBSOCKET] Scheduling reconnection attempt ${reconnectAttempts + 1} in ${delay}ms...`);
    
    reconnectTimeout = setTimeout(() => {
        reconnectTimeout = null;
        reconnectAttempts++;
        connectWebSocket();
    }, delay);
}

function startHTTPPolling() {
    if (pollingTimer) return; // Prevent duplicate timers
    console.log("[POLLING] Active HTTP Polling Fallback Enabled.");
    connDot.className = "w-2.5 h-2.5 bg-amber-500 rounded-full animate-pulse";
    connStatus.textContent = "Connected (HTTP Polling)";

    pollingTimer = setInterval(() => {
        fetch("http://localhost:8080/api/v1/alerts/poll")
            .then(res => {
                if (!res.ok) throw new Error("Gateway offline");
                return res.json();
            })
            .then(data => {
                // Update connection status back to OK if it was offline
                connDot.className = "w-2.5 h-2.5 bg-amber-500 rounded-full animate-pulse";
                connStatus.textContent = "Connected (HTTP Polling)";
                
                if (data && data.event_type) {
                    console.log("[POLLING] Real-Time Alert Received:", data);
                    handleIncomingAlert(data);
                }
            })
            .catch(err => {
                console.log("[POLLING] Gateway unreachable, retrying...");
                connDot.className = "w-2.5 h-2.5 bg-slate-500 rounded-full animate-pulse";
                connStatus.textContent = "Gateway Unreachable";
            });
    }, 1000);
}

function handleIncomingAlert(alert) {
    // Dynamically update the active Room 204 dashboard card
    if (alert.event_type === "alert.fall.predictive" || alert.event_type === "alert.fall.confirmed") {
        setFallState({ bpm: 110, breaths: 26 });
    } else if (alert.event_type === "alert.bed_exit.dangling") {
        setDanglingState({ bpm: 82, breaths: 19 });
    }
}

// Initialize active dashboard states and connect to real-time stream
setStandardState();
connectWebSocket();
