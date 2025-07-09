document.addEventListener("DOMContentLoaded", async () => {
    const container = document.getElementById("shared-notes-container");

    try {
        const res = await fetch("/api/shared", { credentials: "include" });
        const notes = await res.json();

        if (!res.ok) {
            container.textContent = notes.error || "Failed to load shared notes.";
            return;
        }

        if (notes.length === 0) {
            container.innerHTML = `<p>No notes shared with you yet.</p>`;
            return;
        }

        notes.forEach(note => {
            const div = document.createElement("div");
            div.className = "note-card";
            div.innerHTML = `
        <h3>${note.title}</h3>
        <p>${note.content}</p>
        <p class="meta">Shared by <strong>${note.author}</strong> on ${new Date(note.created_at).toLocaleString()}</p>
      `;
            container.appendChild(div);
        });
    } catch (err) {
        container.textContent = "Error loading notes.";
        console.error(err);
    }
});

const shareModal = document.getElementById("share-modal");
const closeShareModalBtn = document.getElementById("close-share-modal");
const userList = document.getElementById("user-list");

let currentNoteToShare = null;

function openShareModal(note) {
    currentNoteToShare = note;
    userList.innerHTML = "<li>Loading users...</li>";
    shareModal.style.display = "flex";

    fetch("/api/users", { credentials: "include" })
        .then(res => res.json())
        .then(users => {
            if (!Array.isArray(users)) {
                userList.innerHTML = "<li>Failed to load users</li>";
                return;
            }
            if (users.length === 0) {
                userList.innerHTML = "<li>No other users found</li>";
                return;
            }
            userList.innerHTML = "";
            users.forEach(user => {
                const li = document.createElement("li");
                li.textContent = `${user.fullname} (${user.username})`;
                li.addEventListener("click", () => {
                    shareNoteWithUser(currentNoteToShare.id, user.id);
                });
                userList.appendChild(li);
            });
        })
        .catch(() => {
            userList.innerHTML = "<li>Error loading users</li>";
        });
}

function closeShareModal() {
    shareModal.style.display = "none";
    currentNoteToShare = null;
}

closeShareModalBtn.addEventListener("click", closeShareModal);

// Clicking outside modal content closes modal
shareModal.addEventListener("click", e => {
    if (e.target === shareModal) {
        closeShareModal();
    }
});

async function shareNoteWithUser(noteId, userId) {
    try {
        const res = await fetch("/api/share-note", {
            method: "POST",
            credentials: "include",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ note_id: noteId, target_user_id: userId }),
        });
        const data = await res.json();
        if (res.ok) {
            alert("Note shared successfully!");
            closeShareModal();
        } else {
            alert(data.error || "Failed to share note");
        }
    } catch (err) {
        alert("Error sharing note");
    }
}

// Attach event listeners to share buttons when rendering notes:

function attachShareButtonListeners() {
    document.querySelectorAll(".share-note-btn").forEach(btn => {
        btn.addEventListener("click", () => {
            // Get the note data related to this button
            const noteCard = btn.closest(".note-card");
            if (!noteCard) return;

            // Assuming you store note data on element dataset (you can adjust accordingly)
            const note = {
                id: Number(noteCard.getAttribute("data-note-id")),
                title: noteCard.querySelector("h3").textContent,
                content: noteCard.querySelector("p").textContent,
            };
            openShareModal(note);
        });
    });
}
