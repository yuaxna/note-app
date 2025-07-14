document.addEventListener("DOMContentLoaded", () => {
  const form = document.getElementById("note-form");
  const titleInput = document.getElementById("note-title");
  const contentInput = document.getElementById("note-content");
  const noteIdInput = document.getElementById("note-id");
  const notesContainer = document.getElementById("notes-container");
  const messageDiv = document.getElementById("note-message");
  const notesCount = document.getElementById("notes-count");

  function showMessage(message, isSuccess = false) {
    messageDiv.textContent = message;
    messageDiv.style.color = isSuccess ? "green" : "red";
    messageDiv.style.display = "block";
    setTimeout(() => {
      messageDiv.textContent = "";
      messageDiv.style.display = "none";
    }, 3000);
  }

  // Export this function so WebSocket can use it
  async function fetchNotes() {
    try {
      const res = await fetch("/api/notes", { credentials: "include" });

      if (!res.ok) {
        showMessage("Failed to fetch notes.");
        return;
      }

      const contentType = res.headers.get("content-type") || "";
      if (!contentType.includes("application/json")) {
        showMessage("Invalid response from server.");
        return;
      }

      const notes = await res.json();
      notesContainer.innerHTML = "";
      if (notesCount) notesCount.textContent = notes.length;

      if (notes.length === 0) {
        notesContainer.innerHTML = `
          <div class="empty-state">
            <div class="empty-icon">📝</div>
            <p>No notes yet</p>
            <p style="font-size: 0.9rem; color: #8B9A7A;">Create your first note using the form on the left</p>
          </div>
        `;
        return;
      }

      notes.forEach(note => {
        const div = document.createElement("div");
        div.className = "note-card";
        div.setAttribute("data-note-id", note.id);

        div.innerHTML = `
          <h3>${note.title} ${!note.is_owner ? '<span class="shared-badge">Shared</span>' : ''}</h3>
          <p>${note.content}</p>
          <p class="note-meta">By ${note.username || "Unknown"} on ${note.created_at || "Unknown date"}</p>
        `;

        if (note.can_edit) {
          const editBtn = document.createElement("button");
          editBtn.textContent = "Edit";
          editBtn.className = "edit-note-btn";
          editBtn.addEventListener("click", () => editNote(note));
          div.appendChild(editBtn);
        }

        if (note.is_owner) {
          const deleteBtn = document.createElement("button");
          deleteBtn.textContent = "Delete";
          deleteBtn.className = "delete-note-btn";
          deleteBtn.addEventListener("click", () => deleteNote(note.id));
          div.appendChild(deleteBtn);

          const shareBtn = document.createElement("button");
          shareBtn.textContent = "Share";
          shareBtn.className = "share-note-btn";
          shareBtn.addEventListener("click", () => {
            if (window.openShareModal) {
              window.openShareModal(note);
            }
          });
          div.appendChild(shareBtn);
        }

        notesContainer.appendChild(div);
      });
    } catch (error) {
      console.error("Error fetching notes:", error);
      showMessage("Error loading notes.");
    }
  }

  function editNote(note) {
    titleInput.value = note.title;
    contentInput.value = note.content;
    noteIdInput.value = note.id;
  }

  async function deleteNote(id) {
    try {
      const res = await fetch(`/api/notes/${id}`, {
        method: "DELETE",
        credentials: "include"
      });

      const contentType = res.headers.get("content-type") || "";
      let data = {};
      if (contentType.includes("application/json")) {
        data = await res.json();
      } else {
        const text = await res.text();
        console.error("Non-JSON response on delete:", text);
      }

      if (res.ok) {
        showMessage("Note deleted", true);
        fetchNotes();
      } else {
        showMessage(data.error || "Failed to delete note");
      }
    } catch (error) {
      console.error("Error deleting note:", error);
      showMessage("Error deleting note");
    }
  }

  form.addEventListener("submit", async (e) => {
    e.preventDefault();

    const noteId = noteIdInput.value;
    const url = noteId ? `/api/notes/${noteId}` : "/api/notes";
    const method = noteId ? "PUT" : "POST";

    const newNote = {
      title: titleInput.value.trim(),
      content: contentInput.value.trim(),
    };

    if (!newNote.title || !newNote.content) {
      showMessage("Please fill in both title and content.");
      return;
    }

    try {
      const res = await fetch(url, {
        method,
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify(newNote),
      });

      const contentType = res.headers.get("content-type") || "";
      let data = {};
      if (contentType.includes("application/json")) {
        data = await res.json();
      } else {
        const text = await res.text();
        console.error("Non-JSON response on save:", text);
        showMessage("Failed to save note: Server returned an invalid response");
        return;
      }

      if (res.ok) {
        titleInput.value = "";
        contentInput.value = "";
        noteIdInput.value = "";
        fetchNotes();
        showMessage(noteId ? "Note updated!" : "Note created!", true);
      } else {
        showMessage(data.error || "Failed to save note");
      }
    } catch (error) {
      console.error("Error saving note:", error);
      showMessage("Error saving note");
    }
  });

  // Initial load
  fetchNotes();

  // Make fetchNotes available globally for WebSocket
  window.fetchNotes = fetchNotes;
});
