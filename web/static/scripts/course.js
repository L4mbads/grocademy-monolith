document.addEventListener("DOMContentLoaded", async () => {
    const paths = window.location.pathname
    .split('/')
    .filter(path => path !== ''); // Remove empty strings from potential trailing slashes

    const lastSegment = paths[paths.length - 1]; // Get the last element
    const courseApiUrl = `/api/courses/${lastSegment}`;

    try {
      const res = await fetch(courseApiUrl);
      const data = await res.json();
      const course = data.data;

      // Fill details
      document.getElementById("thumbnail").src = course.thumbnail_image || "https://upload.wikimedia.org/wikipedia/commons/thumb/f/f5/No-Image-Placeholder-landscape.svg/768px-No-Image-Placeholder-landscape.svg.png";
      document.getElementById("title").textContent = course.title;
      document.getElementById("description").textContent = course.description;
      document.getElementById("instructor").textContent = course.instructor;
      document.getElementById("topics").textContent = course.topics.join(", ");
      document.getElementById("price").textContent = course.price;

      const actionButton = document.getElementById("actionButton");
      const messageEl = document.getElementById("message");

      if (course.purchased) {
        actionButton.textContent = "Open";
        actionButton.onclick = () => {
          window.location.href = `/courses/${course.id}/modules`;
        };
      } else {
        actionButton.textContent = "Buy";
        actionButton.onclick = async () => {
          try {
            const buyRes = await fetch(`/api/courses/${course.id}/buy`, {
              method: "POST",
              headers: { "Content-Type": "application/json" }
            });
            const buyData = await buyRes.json();
            messageEl.textContent = buyData.message;

            if (buyRes.ok) {
              // Optionally redirect to modules after success
              setTimeout(() => {
                window.location.href = `/courses/${course.id}/modules`;
              }, 1000);
            }
          } catch (err) {
            messageEl.textContent = "Error purchasing course.";
          }
        };
      }
    } catch (err) {
      document.getElementById("message").textContent = "Error loading course details.";
    }
  });
