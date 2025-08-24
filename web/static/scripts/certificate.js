var canDownload = false
var completionDate = new Date()

function generateCertificate(username, courseTitle, instructor, date) {
    const canvas = document.createElement("canvas");
    const ctx = canvas.getContext("2d");

    // (A4 ratio: 3508x2480px @ 300dpi, but smaller for web)
    canvas.width = 1200;
    canvas.height = 850;

    // BG
    ctx.fillStyle = "#fdfdfd";
    ctx.fillRect(0, 0, canvas.width, canvas.height);

    // Border
    ctx.strokeStyle = "#000";
    ctx.lineWidth = 10;
    ctx.strokeRect(20, 20, canvas.width - 40, canvas.height - 40);

    // Title
    ctx.fillStyle = "#333";
    ctx.font = "bold 48px Serif";
    ctx.textAlign = "center";
    ctx.fillText("Certificate of Completion", canvas.width / 2, 120);

    // Subtitle
    ctx.font = "24px Serif";
    ctx.fillText("This is proudly presented to", canvas.width / 2, 200);

    // Username
    ctx.font = "bold 40px Serif";
    ctx.fillText(username, canvas.width / 2, 280);

    // Course info
    ctx.font = "22px Serif";
    ctx.fillText(`For successfully completing the course`, canvas.width / 2, 350);

    ctx.font = "italic 28px Serif";
    ctx.fillText(courseTitle, canvas.width / 2, 400);

    // Instructor and Date
    ctx.font = "20px Serif";
    ctx.textAlign = "left";
    ctx.fillText(`Instructor: ${instructor}`, 100, canvas.height - 120);
    ctx.fillText(`Date: ${date}`, 100, canvas.height - 80);

    ctx.textAlign = "right";
    ctx.fillText("Signed", canvas.width - 100, canvas.height - 80);

    return canvas;
}

function downloadCertificate(username, courseTitle, instructor, date) {
    const canvas = generateCertificate(username, courseTitle, instructor, date);
    const link = document.createElement("a");
    link.download = `certificate_${courseTitle.replace(/\s+/g, "_")}.png`;
    link.href = canvas.toDataURL("image/png");
    link.click();
}

async function handleDownloadCertificate() {
    if (!canDownload) return;
    console.log(completionDate)
    const user = await fetchUserData();
    downloadCertificate(user.username, document.getElementById("course-title").textContent, document.getElementById("course-instructor").textContent.split(' ').slice(1).join(' '), completionDate.toLocaleDateString())
}

function setDownload(val) {
    canDownload = val
    document.getElementById("downloadCertificate").disabled = !val;
}