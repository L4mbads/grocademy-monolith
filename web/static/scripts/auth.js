async function handleAuthForm(event, endpoint, messageId) {
    event.preventDefault();

    const form = event.target;
    const messageEl = document.getElementById(messageId);
    messageEl.textContent = "";
    messageEl.className = "form-message"; // reset

    const formData = new FormData(form);
    const data = Object.fromEntries(formData.entries());

    if (document.querySelector(".register")) {
      if(!isStrongPassword(data.password)) {
          messageEl.textContent = "Password must be at least 8 characters containing number, lowercase, and uppercase letters.";
          messageEl.classList.add("error");
          return
      }

      if(data.password !== data.confirm_password) {
          messageEl.textContent = "Confirm password doesn't match.";
          messageEl.classList.add("error");
          return
      }
    }

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


async function fetchUserData() {
  try {
    const res = await fetch("/api/auth/self");
    const result = await res.json();

    if (!res.ok) {
      console.error("Error fetching user:", result.message);
      return {};
    }

    return result.data;
  } catch (err) {
    console.error("Failed to fetch user data", err);
  }
}

function isStrongPassword(password) {
  const hasMinLen = password.length >= 8;
  const hasUpper = /[A-Z]/.test(password);
  const hasLower = /[a-z]/.test(password);
  const hasNumber = /[0-9]/.test(password);

  return hasMinLen && hasUpper && hasLower && hasNumber;
}