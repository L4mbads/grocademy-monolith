async function queryModule() {
    const courseId = getCourseIdFromUrl();
    const container = document.getElementById("modules-container");
    const courseTitle = document.getElementById("course-title");
    const progressFill = document.getElementById("progress-fill");
    const progressText = document.getElementById("progress-text");

    const route = `api/courses/${courseId}/modules`;
    const body = await query(route)

    courseTitle.textContent = body.message.split(" ")[0];
    container.innerHTML = "";

    if (body.status !== "success" || !body.data) {
        container.innerHTML = `<p class="form-message error">Failed to load modules: ${body.message}</p>`;
        return;
    }

    const modules = body.data.sort((a, b) => a.order - b.order);

    const totalModules = body.data.length
    var completedModules = 0

    modules.forEach((mod) => {

        if (mod.is_completed) {
            completedModules += 1;
        }

        const card = document.createElement("div");
        card.className = "module-card"

        const title = document.createElement("h2");
        title.textContent = mod.title;
        title.className = "module-title"
        card.appendChild(title)

        const description = document.createElement("p");
        description.textContent = mod.description;
        description.className = "module-description"
        card.appendChild(description)

        const actions = document.createElement("div");
        actions.className = "module-actions"
        card.appendChild(actions)

        const pdf = document.createElement("a");
        pdf.href = mod.pdf_content;
        pdf.download = "pdf";
        pdf.className = "btn"
        pdf.innerText = "Download PDF"
        actions.appendChild(pdf)

        const video = document.createElement("a");
        video.href = mod.video_content;
        video.download = "video";
        video.className = "btn"
        video.innerText = "Download Video"
        actions.appendChild(video)

        const actionButton = document.createElement("button")
        actionButton.className = "btn complete-btn";
        actionButton.innerText = mod.is_completed ? "Completed" : "Mark Complete";
        actionButton.dataset.id = mod.id
        actions.appendChild(actionButton)

        container.appendChild(card);
    });

    const percentage = completedModules / totalModules;
    progressFill.style = `width: ${percentage*100}%;`
    progressText.textContent = `${percentage*100}%`

    // Event listener for marking module complete
    container.addEventListener("click", async (e) => {
        if (e.target.classList.contains("complete-btn")) {
            const moduleId = e.target.dataset.id;

            try {
                const res = await fetch(`/api/modules/${moduleId}/complete`, {
                    method: "PATCH",
                    body: JSON.stringify({"is_completed": e.target.textContent !== "Completed"}),
                });
                const result = await res.json();

                if (result.status === "success") {
                    e.target.textContent = result.data.is_completed ? "Completed" : "Mark Complete";
                    // e.target.disabled = true;

                    progressFill.style = `width: ${result.data.course_progress.percentage}%;`
                    progressText.textContent = `${result.data.course_progress.percentage}%`



                } else {
                    alert("Failed: " + result.message);
                }
            } catch (err) {
                console.error(err);
                alert("Error completing module");
            }
        }
    });
}

// function startLongPolling() {
//     fetch('/api/updates') // Endpoint for long polling on the server
//         .then(response => response.json())
//         .then(data => {
//         // Process the received data and update the UI
//         console.log('New data received:', data);
//         // Immediately initiate a new long polling request
//         startLongPolling();
//         })
//         .catch(error => {
//         console.error('Long polling error:', error);
//         // Implement retry logic or backoff strategies if needed
//         setTimeout(startLongPolling, 5000); // Retry after 5 seconds
//         });
//     }

// document.addEventListener('DOMContentLoaded', startLongPolling);
document.addEventListener('DOMContentLoaded', queryModule);

function handleSearchInput(e){
    if(e.keyCode === 13){
        e.preventDefault();
        search();
        queryModule();
    }
}

function getCourseIdFromUrl() {
    const parts = window.location.pathname.split("/");
    return parts[2];
}

function handleSearchInput(e){
    if(e.keyCode === 13){
        e.preventDefault();
        handleSearch();
    }
}

function handleSearch() {
    search();
    queryModule();
}

function handlePrevBtn() {
    changePage(-1);
    queryModule();
}

function handleNextBtn() {
    changePage(1);
    queryModule();
}

function handleLimitChange() {
    changeLimit();
    queryModule();
}