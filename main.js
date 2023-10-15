document.addEventListener("DOMContentLoaded", function() {
    const uploadForm = document.getElementById("uploadForm");
    const mp3FileInput = document.getElementById("mp3File");
    const audio = document.getElementById("audio");
    const startTimeInput = document.getElementById("startTime");
    const endTimeInput = document.getElementById("endTime");
    const slowFilterButton = document.getElementById("slowFilter");
    let selectedFile;
    let startTime = 0;
    let endTime = 30;
    mp3FileInput.addEventListener("change", function(event) {
        event.preventDefault();
        selectedFile = mp3FileInput.files[0];
        audio.src = URL.createObjectURL(selectedFile);
        audio.currentTime = startTime;
    });
    startTimeInput.addEventListener("input", function() {
        startTime = parseFloat(startTimeInput.value);
        endTime = startTime + 30;
        endTimeInput.value = endTime;
        audio.currentTime = startTime;
    });
    audio.addEventListener("timeupdate", function() {
        if (audio.currentTime < startTime) { audio.currentTime = startTime }
        if (audio.currentTime > endTime) { audio.currentTime = startTime; audio.pause() }
    });
    slowFilterButton.addEventListener("click", async function() {
        if (!selectedFile) {
            console.log("No file selected");
            return;
        }
        const formData = new FormData();
        formData.append("file", selectedFile, "clipped_audio.mp3");
        formData.append("startTime", startTime.toString());
        formData.append("endTime", endTime.toString());
        formData.append("Content-Type", "audio/mpeg");
        try {
            const response = await fetch("http://localhost:8080/edit_track/0", {
                method: "POST",
                body: formData,
            });
            if (response.ok) {
                console.log("Edit successful");
                // Print details about the processed file
                const contentDisposition = response.headers.get("Content-Disposition");
                console.log("Content-Disposition:", contentDisposition);
                const contentLength = response.headers.get("Content-Length");
                console.log("Content-Length:", contentLength);
                const contentType = response.headers.get("Content-Type");
                console.log("Content-Type:", contentType);
            } else {
                console.log("Edit failed");
                console.log(response);
            }
        } catch (err) {
            console.log("Error: ", err);
        }
    });
});
