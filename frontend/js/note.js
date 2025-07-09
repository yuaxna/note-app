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
        div.setAttribute("data-note-id", note.id);

        div.innerHTML = `
          <h3>${note.title}</h3>
          <p>${note.content}</p>
          <p class="note-meta">By ${note.username || "Unknown"} on ${note.created_at || "Unknown date"}</p>
          <button class="edit-note-btn">Edit</button>
          <button class="delete-note-btn">Delete</button>
          <button class="share-note-btn">Share</button>
        `;

        div.querySelector(".delete-note-btn").addEventListener("click", () => {
          deleteNote(note.id);
        });

        div.querySelector(".edit-note-btn").addEventListener("click", () => {
          editNote(note);
        });

        div.querySelector(".share-note-btn").addEventListener("click", () => {
          const identifier = prompt("Enter the username or email of the person to share with:");
          if (!identifier) return;

          fetch("/api/notes/share", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify({ note_id: note.id, identifier }),
          })
          .then(res => res.json())
          .then(data => {
            if (data.message) {
              alert(data.message);
            } else {
              alert(data.error || "Failed to share note.");
            }
          })
          .catch(err => {
            console.error("Share error:", err);
            alert("An error occurred while sharing the note.");
          });
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

  // Sidebar toggle
  const menuBtn = document.getElementById("menu-btn");
  const sidebar = document.getElementById("sidebar");

  if (menuBtn && sidebar) {
    menuBtn.addEventListener("click", () => {
      sidebar.classList.toggle("sidebar-closed");
    });
  }
});
