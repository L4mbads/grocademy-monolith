document.addEventListener("DOMContentLoaded", async () => {
    const toggleBtn = document.getElementById("nav-toggle");
    const navLinks = document.getElementById("nav-links");

    if (toggleBtn) {
      toggleBtn.addEventListener("click", () => {
        navLinks.classList.toggle("show");
      });
    }

    const res = await fetch("/api/auth/self");
    const result = await res.json();
    if (!res.ok) {
      console.error("Error fetching user:", result.message);
    } else {
      const user = result.data;
      document.getElementById("userdata").innerHTML = `${user.username}: $${user.balance}`;
    }
});