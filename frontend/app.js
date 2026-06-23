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
function setDanglingState() {
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
    vitalHR.textContent = "82 BPM";
    vitalRR.textContent = "19 /min";
}

// Set Red Danger state for Fall
function setFallState() {
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
    vitalHR.textContent = "110 BPM";
    vitalRR.textContent = "26 /min";
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
document.getElementById('btn-standard').addEventListener('click', setStandardState);
document.getElementById('btn-dangling').addEventListener('click', setDanglingState);
document.getElementById('btn-fall').addEventListener('click', setFallState);
document.getElementById('btn-ack').addEventListener('click', setAcknowledgedState);

// Default boot setup
setStandardState();
