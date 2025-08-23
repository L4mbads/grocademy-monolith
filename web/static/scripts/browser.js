let currentPage = 1;
let totalPages = 1;
let currentQuery = "";
let currentLimit = 15;

async function query(route) {
    const res = await fetch(`${encodeURIComponent(route)}?q=${encodeURIComponent(currentQuery)}&page=${currentPage}&limit=${currentLimit}`);
    const data = await res.json();

    return data
}

function search() {
    currentQuery = document.getElementById("searchInput").value.trim();
    currentPage = 1;
}

function changePage(delta) {
    if ((currentPage + delta) >= 1 && (currentPage + delta) <= totalPages) {
        currentPage += delta;
    }
}

function changeLimit() {
    currentLimit = parseInt(document.getElementById("limitSelect").value);
    currentPage = 1;
}