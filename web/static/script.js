async function handleAuthForm(event, endpoint, messageId) {
    event.preventDefault();

    const form = event.target;
    const messageEl = document.getElementById(messageId);
    messageEl.textContent = "";
    messageEl.className = "form-message"; // reset

    const formData = new FormData(form);
    const data = Object.fromEntries(formData.entries());

    try {
      const res = await fetch(endpoint, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      });

      const result = await res.json();

      if (!res.ok) {
        messageEl.textContent = result.message || "Something went wrong.";
        messageEl.classList.add("error");
        return;
      }

      messageEl.textContent = result.message || "Success!";
      messageEl.classList.add("success");

      // Example: redirect after short delay
      setTimeout(() => {
        if (endpoint.includes("login")) {
          window.location.href = "/dashboard";
        } else {
          window.location.href = "/login";
        }
      }, 1200);

    } catch (err) {
      console.error("Request failed:", err);
      messageEl.textContent = "Network error.";
      messageEl.classList.add("error");
    }
  }

  document.addEventListener("DOMContentLoaded", () => {
    const loginForm = document.querySelector("#login-form");
    const registerForm = document.querySelector("#register-form");

    if (loginForm) {
      loginForm.addEventListener("submit", (e) =>
        handleAuthForm(e, "/api/auth/login", "login-message")
      );
    }
    if (registerForm) {
      registerForm.addEventListener("submit", (e) =>
        handleAuthForm(e, "/api/auth/register", "register-message")
      );
    }
    const menu_toggle = document.querySelector('.menu-toggle');
    const sidebar = document.querySelector('.sidebar');
    if (menu_toggle) {
        menu_toggle.addEventListener('click', () => {
            menu_toggle.classList.toggle('is-active');
            sidebar.classList.toggle('is-active');
        });
    }
  });

  // Toggle Navbar
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

  async function fetchUserData() {
    try {
      const res = await fetch("/api/auth/self");
      const result = await res.json();

      if (!res.ok) {
        console.error("Error fetching user:", result.message);
        return;
      }

      const user = result.data;

      document.getElementById("user-name").textContent = user.first_name || user.username;
      document.getElementById("user-username").textContent = user.username;
      document.getElementById("user-email").textContent = user.email;
      document.getElementById("user-firstname").textContent = user.first_name;
      document.getElementById("user-lastname").textContent = user.last_name;
      document.getElementById("user-balance").textContent = user.balance;
      document.getElementById("user-created").textContent = new Date(user.created_at).toLocaleDateString();
    } catch (err) {
      console.error("Failed to fetch user data", err);
    }
  }