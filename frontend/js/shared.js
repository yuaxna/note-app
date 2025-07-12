const shareModal = document.getElementById("share-modal");
const closeShareModalBtn = document.getElementById("close-share-modal");
const userList = document.getElementById("user-list");

if (shareModal && closeShareModalBtn && userList) {
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

  shareModal.addEventListener("click", e => {
    if (e.target === shareModal) {
      closeShareModal();
    }
  });

  async function shareNoteWithUser(noteId, userId) {
    try {
      const res = await fetch("/api/share", {
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

  async function loadSharedNotes() {
    const container = document.getElementById("shared-notes-container");
    if (!container) return;

    container.innerHTML = "Loading shared notes...";

    try {
      const res = await fetch("/api/shared", { credentials: "include" });
      if (!res.ok) throw new Error("Failed to load shared notes");

      const notes = await res.json();

      if (!notes.length) {
        container.innerHTML = "<p>No shared notes yet.</p>";
        return;
      }

      container.innerHTML = "";
      notes.forEach(note => {
        const noteDiv = document.createElement("div");
        noteDiv.classList.add("note-card");
        noteDiv.setAttribute("data-note-id", note.id);
        noteDiv.innerHTML = `
          <h3>${note.title}</h3>
          <p>${note.content}</p>
          <small>By: ${note.author} on ${new Date(note.created_at).toLocaleString()}</small>
          <button class="share-note-btn">Share</button>
        `;
        container.appendChild(noteDiv);
      });
    } catch (err) {
      container.innerHTML = `<p>Error loading shared notes: ${err.message}</p>`;
    }
  }

  document.addEventListener("DOMContentLoaded", () => {
    loadSharedNotes();
  });

  // 🔓 Make openShareModal globally accessible to note.js
  window.openShareModal = openShareModal;
}
