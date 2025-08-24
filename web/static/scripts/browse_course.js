
async function queryCourse() {
    console.log(route)
    const data = await query(route)
    const container = document.getElementById("coursesContainer");
    container.innerHTML = "";

    if (data.status === "success" && data.data) {
        data.data.forEach(course => {
        const card = document.createElement("div");
        card.className = "card";

        const img = document.createElement("img");
        img.src = course.thumbnail_image || "https://upload.wikimedia.org/wikipedia/commons/thumb/f/f5/No-Image-Placeholder-landscape.svg/768px-No-Image-Placeholder-landscape.svg.png";
        card.appendChild(img);

        const cardDetail = document.createElement("div");
        cardDetail.className = "card-detail";
        card.appendChild(cardDetail)

        const title = document.createElement("h3");
        title.textContent = course.title;
        cardDetail.appendChild(title);

        const instructor = document.createElement("p");
        instructor.textContent = "Instructor: " + course.instructor;
        cardDetail.appendChild(instructor);

        const desc = document.createElement("p");
        desc.textContent = course.description;
        cardDetail.appendChild(desc);

        const price = document.createElement("div");
        price.className = "price";
        price.textContent = course.price > 0 ? "$" + course.price : "Free";
        cardDetail.appendChild(price);

        const btn = document.createElement("button");
        btn.textContent = "View Details";
        btn.onclick = () => {
            window.location.href = `/courses/${course.id}`;
        };

        cardDetail.appendChild(btn);

        container.appendChild(card);
        });

        // update pagination
        totalPages = data.pagination.total_pages;
        document.getElementById("pageInfo").textContent = `Page ${data.pagination.current_page} of ${totalPages}`;
        document.getElementById("prevBtn").disabled = currentPage === 1;
        document.getElementById("nextBtn").disabled = currentPage === totalPages;
    } else {
        container.innerHTML = `<p>${"No courses found"}</p>`;
    }

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
document.addEventListener('DOMContentLoaded', queryCourse);

function handleSearchInput(e){
    if(e.keyCode === 13){
        e.preventDefault();
        handleSearch();
    }
}

function handleSearch() {
    search();
    queryCourse();
}

function handlePrevBtn() {
    changePage(-1);
    queryCourse();
}

function handleNextBtn() {
    changePage(1);
    queryCourse();
}

function handleLimitChange() {
    changeLimit();
    queryCourse();
}