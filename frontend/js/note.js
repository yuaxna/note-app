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

  async function fetchNotes() {
    try {
      const res = await fetch("/api/notes", {
        method: "GET",
        credentials: "include",
      });

      if (!res.ok) {
        showMessage("Failed to fetch notes.");
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

        const date = note.created_at ? new Date(note.created_at) : null;
        const formattedDate = date ? date.toLocaleString() : "Invalid Date";
        const author = note.username || "Unknown";

        div.innerHTML = `
          <h3>${note.title}</h3>
          <p>${note.content}</p>
          <div class="note-meta">
            <small>By <strong>${author}</strong> on <em>${formattedDate}</em></small>
          </div>
          <button class="edit-note-btn">Edit</button>
          <button class="delete-note-btn">Delete</button>
        `;

        // Event listeners
        div.querySelector(".delete-note-btn").addEventListener("click", () => {
          deleteNote(note.id);
        });

        div.querySelector(".edit-note-btn").addEventListener("click", () => {
          editNote(note);
        });

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

      const data = await res.json();

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
    const method = noteId ? "PUT" : "POST";
    const url = "/api/notes";

    const newNote = {
      title: titleInput.value.trim(),
      content: contentInput.value.trim(),
    };
    if (noteId) newNote.id = parseInt(noteId);

    if (!newNote.title || !newNote.content) {
      showMessage("Please fill in both title and content.");
      return;
    }

    try {
      const res = await fetch(url, {
        method,
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify(newNote),
      });

      const data = await res.json();

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

  // Initial fetch
  fetchNotes();
});

// Sidebar toggle
const menuBtn = document.getElementById("menu-btn");
const sidebar = document.getElementById("sidebar");

if (menuBtn && sidebar) {
  menuBtn.addEventListener("click", () => {
    sidebar.classList.toggle("sidebar-closed");
  });
}
