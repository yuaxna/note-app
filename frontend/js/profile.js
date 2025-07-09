document.addEventListener('DOMContentLoaded', async () => {
    const menuBtn = document.getElementById('menu-btn');
    const sidebar = document.getElementById('sidebar');

    // Sidebar toggle
    menuBtn.addEventListener('click', () => {
        sidebar.classList.toggle('sidebar-closed');
    });

    // Logout functionality
    document.getElementById("logout-btn").addEventListener("click", async () => {
        const res = await fetch("/logout", {
            method: "POST",
            credentials: "include"
        });

        const data = await res.json();

        if (res.ok) {
            alert("Logged out successfully.");
            window.location.href = "/";
        } else {
            alert(data.error || "Logout failed.");
        }
    });

    // Load user info
    try {
        const res = await fetch("/api/me", {
            credentials: "include"
        });

        if (res.ok) {
            const user = await res.json();
            document.getElementById("username").textContent = user.username || "User";
            document.getElementById("info-fullname").textContent = user.fullname || "N/A";
            document.getElementById("info-username").textContent = user.username || "N/A";
            document.getElementById("info-email").textContent = user.email || "N/A";
            // Add other user info fields as needed
        } else {
            console.error("Failed to load user info");
            // Don't show alert here, just log the error
        }
    } catch (error) {
        console.error("Error loading user info:", error);
    }

    // Load user notes
    try {
        const res = await fetch("/api/notes", {
            credentials: "include"
        });

        if (res.ok) {
            const notes = await res.json();
            const notesContainer = document.getElementById("user-notes-list");

            if (notes.length === 0) {
                notesContainer.innerHTML = "<p>No notes yet.</p>";
            } else {
                notesContainer.innerHTML = notes.map(note => `
                    <div class="note-item">
                        <h4>${note.title}</h4>
                        <p>${note.content}</p>
                    </div>
                `).join('');
            }
        }
    } catch (error) {
        console.error("Error loading notes:", error);
    }
});
