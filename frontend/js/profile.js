document.addEventListener("DOMContentLoaded", async () => {
    try {
        // Get user info
        const meRes = await fetch("/api/me");
        const me = await meRes.json();

        if (!me.name) {
            alert("You are not logged in.");
            window.location.href = "/";
            return;
        }

        // Fill user info
        document.getElementById("info-fullname").textContent = me.name;
        document.getElementById("info-username").textContent = me.username;
        document.getElementById("info-email").textContent = me.email;
        document.getElementById("info-gender").textContent = me.gender;
        document.getElementById("info-joined").textContent = me.joined;

        // Get notes
        const notesRes = await fetch(`/api/notes?user_id=${me.id}`);
        const notes = await notesRes.json();

        const container = document.getElementById("user-notes-list");
        notes.forEach(note => {
            const div = document.createElement("div");
            div.className = "note-item";
            div.innerHTML = `<strong>${note.title}</strong><p>${note.content}</p>`;
            container.appendChild(div);
        });
    } catch (err) {
        console.error("Error loading profile:", err);
    }
});
