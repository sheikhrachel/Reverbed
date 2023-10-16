const numberOfOrbs = 10;
const radius = 20; // px
const delayIncrement = 0.05; // seconds
const EditTypes = {
    slowed: 0,
    spedUp: 1,
    reverb: 2,
    reversed: 3,
    pitchedUp: 4,
    pitchedDown: 5
};
let deg = 0;

document.addEventListener("DOMContentLoaded", function() {
    let selectedFile;
    const spinner = document.getElementById('loading-spinner');
    // Create orbs
    for (let i = 0; i < numberOfOrbs; i++) {
        const orb = document.createElement('div');
        orb.className = 'orb';
        orb.style.animation = `pulse 1s ease-in-out ${i * delayIncrement}s infinite alternate`;
        spinner.appendChild(orb);
    }
    // Position orbs
    const orbs = document.querySelectorAll('.orb');
    orbs.forEach((orb, index) => {
        const angle = (index / numberOfOrbs) * 2 * Math.PI;
        orb.style.left = `${50 + radius * Math.cos(angle)}%`;
        orb.style.top = `${50 + radius * Math.sin(angle)}%`;
    });
    // Rotate spinner
    setInterval(() => { spinner.style.transform = `rotate(${(deg + 1) % 360}deg)` }, 16);
    const mp3FileInput = document.getElementById("file_input");
    const audio = new Audio();
    const audioControls = document.getElementById('audio-controls');
    const buttonContainer = document.getElementById('filter-buttons');
    const audioVisualizer = document.getElementById('audio-visualizer');
    audioControls.style.display = 'none';
    buttonContainer.style.display = 'none';
    audioVisualizer.style.display = 'none';
    mp3FileInput.addEventListener("change", function(event) {
        event.preventDefault();
        selectedFile = mp3FileInput.files[0];
        if (selectedFile) {
            audio.src = URL.createObjectURL(selectedFile);
            audioControls.style.display = 'block';
            buttonContainer.style.display = 'block';
            audioVisualizer.style.display = 'block';
            buttonContainer.innerHTML = '';
            for (const [editTypeKey, editTypeVal] of Object.entries(EditTypes)) {
                buttonContainer.append(create_filter_button(
                    selectedFile,
                    editTypeKey,
                    audio,
                    editTypeVal,

                ));
            }
        }
    });
    const audioContext = new (window.AudioContext || window.webkitAudioContext)();
    const analyser = audioContext.createAnalyser();
    const audioSrc = audioContext.createMediaElementSource(audio);
    const canvas = document.getElementById("visualizerCanvas");
    const canvasCtx = canvas.getContext("2d");
    audioSrc.connect(analyser);
    analyser.connect(audioContext.destination);
    function draw() {
        // Get the frequency data
        const bufferLength = analyser.frequencyBinCount;
        const dataArray = new Uint8Array(bufferLength);
        analyser.getByteFrequencyData(dataArray);
        // Clear the canvas
        canvasCtx.clearRect(0, 0, canvas.width, canvas.height);
        // Draw bars
        const barWidth = (canvas.width / bufferLength) * 10;
        let barHeight;
        let x = 0;
        for(let i = 0; i < bufferLength; i++) {
            barHeight = dataArray[i];
            const gradient = canvasCtx.createLinearGradient(0, 0, 0, canvas.height);
            gradient.addColorStop(0, "#c091eb");
            gradient.addColorStop(1, "#659ba6");
            canvasCtx.fillStyle = gradient;
            canvasCtx.fillRect(x, canvas.height - barHeight / 2, barWidth, barHeight / 2);
            x += barWidth + 1;
        }
        requestAnimationFrame(draw);
    }
    audio.addEventListener('play', () => {
        audioContext.resume().then(() => {
            draw();
        });
    });
    const playPauseButton = document.getElementById('playPauseButton');
    const progressBar = document.getElementById('progressBar');
    const progressContainer = document.getElementById('progressContainer');
    playPauseButton.addEventListener('click', () => {
        if (!selectedFile) {
            console.log("No file selected");
            return;
        }
        if (audio.paused) {
            audio.play().catch(err => console.error("Playback failed:", err));
            playPauseButton.textContent = 'Pause';
        } else {
            audio.pause();
            playPauseButton.textContent = 'Play';
        }
    });
    audio.addEventListener('timeupdate', () => {
        const progress = (audio.currentTime / audio.duration) * 100;
        progressBar.style.width = progress + '%';
    });
    progressContainer.addEventListener('click', (e) => {
        const clickPosition = e.clientX - progressContainer.getBoundingClientRect().left;
        audio.currentTime = (clickPosition / progressContainer.clientWidth) * audio.duration;
    });
});

function create_filter_button(selectedFile, filter, audio, editTypeVal) {
    const button = document.createElement('button');
    button.className = 'filter-btn filter-btn-blue';
    let readableFilter = filter.replace(/([a-z])([A-Z])/g, '$1 $2'); // Convert camelCase to space-separated words
    readableFilter = readableFilter.toLowerCase();
    button.id = readableFilter;
    button.textContent = readableFilter;
    button.addEventListener("click", async function () {
        document.getElementById("loading-spinner").style.display = "block";
        await send_mp3_file(selectedFile, audio, editTypeVal);
        document.getElementById("loading-spinner").style.display = "none";
    });
    return button
}

async function send_mp3_file(selectedFile, audio, editType) {
    if (!selectedFile) {
        console.log("No file selected");
        return;
    }
    const formData = new FormData();
    formData.append("file", selectedFile, "clipped_audio.mp3");
    formData.append("Content-Type", "audio/mpeg");
    try {
        const response = await fetch("http://localhost:8080/edit_track/" + editType, {
            method: "POST",
            body: formData,
        });
        if (response.ok) {
            const blob = await response.blob();
            audio.src = URL.createObjectURL(blob);
        } else {
            console.log(response);
        }
    } catch (err) {
        console.log("Error: ", err);
    } finally {
        document.getElementById("loading-spinner").style.display = "none";
        audio.pause();
        audio.play().catch(err => console.error("Playback failed:", err));
    }
}
