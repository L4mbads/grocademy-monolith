document.addEventListener("DOMContentLoaded", () => {
    const toggleBtn = document.getElementById("nav-toggle");
    const navLinks = document.getElementById("nav-links");

    if (toggleBtn) {
      toggleBtn.addEventListener("click", () => {
        navLinks.classList.toggle("show");
      });
    }

    // Fetch user data only if dashboard is loaded
    if (document.querySelector(".dashboard")) {
      fetchUserData();
    }
});