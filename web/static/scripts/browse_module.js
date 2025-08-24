
async function queryModule() {
    const courseId = getCourseIdFromUrl();
    const courseTitle = document.getElementById("course-title");
    const courseInstructor = document.getElementById("course-instructor");
    const container = document.getElementById("modules-container");
    const progressFill = document.getElementById("progress-fill");
    const progressText = document.getElementById("progress-text");

    const route = `api/courses/${courseId}/modules`;
    const body = await query(route)


    const res = await fetch(`/api/courses/${courseId}`, {
        method: "GET",
    })
    const course = await res.json()

    courseTitle.textContent = course.data.title;
    courseInstructor.textContent = `by ${course.data.instructor}`;
    container.innerHTML = "";

    if (body.status !== "success" || !body.data) {
        container.innerHTML = `<p class="form-message error">Failed to load modules: ${body.message}</p>`;
        return;
    }

    const modules = body.data.sort((a, b) => a.order - b.order);

    const totalModules = body.data.length
    var completedModules = 0
    var dummyModuleID = -1

    modules.forEach((mod) => {

        if (mod.is_completed) {
            completedModules += 1
            dummyModuleID = mod.id

        }

        const card = document.createElement("div");
        card.className = "module-card"

        const title = document.createElement("h2");
        title.textContent = mod.title;
        title.className = "module-title"
        card.appendChild(title)


        if (mod.video_content) {
            // const video = document.createElement("a");
            // video.href = mod.video_content;
            // video.download = "video";
            // video.className = "btn"
            // video.innerText = "Download Video"
            const videoContainer = document.createElement("div");
            videoContainer.className = "video-container"
            card.appendChild(videoContainer)

            const video = document.createElement("video");
            video.controls = true
            video.textContent = "Your browser does not support video tag."
            videoContainer.appendChild(video)

            const source = document.createElement("source");
            source.src = mod.video_content;
            video.appendChild(source);
        }

        const description = document.createElement("p");
        description.textContent = mod.description;
        description.className = "module-description"
        card.appendChild(description)

        const actions = document.createElement("div");
        actions.className = "module-actions"
        card.appendChild(actions)

        if (mod.pdf_content) {
            const pdf = document.createElement("a");
            pdf.href = mod.pdf_content;
            pdf.className = "btn"
            pdf.target = "_blank"
            pdf.rel = "noopener noreferrer"
            pdf.innerText = "Download PDF"
            actions.appendChild(pdf)
        }


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

    if (percentage === 1 && dummyModuleID >= 0) {
        try {
            const res = await fetch(`/api/modules/${dummyModuleID}/complete`, {
                method: "PATCH",
                body: JSON.stringify({"is_completed": true }),
            });
            const result = await res.json();

            if (result.status === "success") {

                if (result.data.course_progress.percentage === 100) {
                    completionDate = new Date(result.data.course_progress.latest_completion)
                    setDownload(true)
                } else {
                    setDownload(false)
                }

            } else {
                alert("Failed: " + result.message);
            }
        } catch (err) {
            console.error(err);
            alert("Error reading module");
        }
    } else {
        setDownload(false)
    }

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


                    if (result.data.course_progress.percentage === 100) {
                        completionDate = new Date(result.data.course_progress.latest_completion)
                        setDownload(true)
                    } else {
                        setDownload(false)
                    }

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

document.addEventListener('DOMContentLoaded', queryModule);