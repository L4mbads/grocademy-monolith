document.addEventListener("DOMContentLoaded", async () => {
    const toggleBtn = document.getElementById("nav-toggle");
    const navLinks = document.getElementById("nav-links");

    if (toggleBtn) {
      toggleBtn.addEventListener("click", () => {
        navLinks.classList.toggle("show");
      });
    }

    // Fetch user data only if dashboard is loaded
    if (document.querySelector(".dashboard")) {
      const user = await fetchUserData();
      console.log(user)

      document.getElementById("user-name").textContent = user.first_name || user.username;
      document.getElementById("user-username").textContent = user.username;
      document.getElementById("user-email").textContent = user.email;
      document.getElementById("user-firstname").textContent = user.first_name;
      document.getElementById("user-lastname").textContent = user.last_name;
      document.getElementById("user-balance").textContent = user.balance;
      document.getElementById("user-created").textContent = new Date(user.created_at).toLocaleDateString();
    }
});