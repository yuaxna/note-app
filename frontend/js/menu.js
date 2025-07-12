document.addEventListener("DOMContentLoaded", () => {
    // Sidebar toggle
    const menuBtn = document.getElementById("menu-btn");
    const sidebar = document.getElementById("sidebar");

    if (menuBtn && sidebar) {
        menuBtn.addEventListener("click", () => {
            sidebar.classList.toggle("sidebar-closed");
        });
    }

    // Logout
    const logoutBtn = document.getElementById("logout-btn");
    if (logoutBtn) {
        logoutBtn.addEventListener("click", async () => {
            try {
                const res = await fetch("/logout", {
                    method: "POST",
                    credentials: "include",
                });

                if (res.ok) {
                    alert("Logged out successfully.");
                    window.location.href = "/";
                } else {
                    const data = await res.json();
                    alert(data.error || "Logout failed.");
                }
            } catch (err) {
                console.error("Logout error:", err);
                alert("Error during logout.");
            }
        });
    }

    // Load username into sidebar
    const sidebarUsername = document.getElementById("username");
    if (sidebarUsername) {
        fetch("/api/me", {
            credentials: "include"
        })
        .then(res => res.ok ? res.json() : null)
        .then(user => {
            if (user && user.username) {
                sidebarUsername.textContent = user.username;
            }
        })
        .catch(err => {
            console.error("Error loading sidebar username:", err);
        });
    }
});
