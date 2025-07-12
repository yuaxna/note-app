document.addEventListener('DOMContentLoaded', async () => {
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
            document.getElementById("info-gender").textContent = user.gender || "N/A";
        } else {
            console.error("Failed to load user info");
            document.getElementById("info-fullname").textContent = "Error loading";
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
            const notesCount = document.getElementById("notes-count");

            notesCount.textContent = notes.length;

            if (notes.length === 0) {
                notesContainer.innerHTML = `
                            <div class="empty-state">
                                <div class="empty-icon">📝</div>
                                <p>No notes yet</p>
                                <a href="/home" class="create-note-link">Create your first note</a>
                            </div>
                        `;
            } else {
                notesContainer.innerHTML = notes.map(note => `
                            <div class="note-item">
                                <div class="note-header">
                                    <h4>${note.title}</h4>
                                    <span class="note-date">${new Date(note.created_at).toLocaleDateString()}</span>
                                </div>
                                <p class="note-content">${note.content}</p>
                            </div>
                        `).join('');
            }
        } else {
            document.getElementById("user-notes-list").innerHTML = `
                        <div class="error-state">
                            <p>Error loading notes.</p>
                        </div>
                    `;
        }
    } catch (error) {
        console.error("Error loading notes:", error);
        document.getElementById("user-notes-list").innerHTML = `
                    <div class="error-state">
                        <p>Error loading notes.</p>
                    </div>
                `;
    }
});