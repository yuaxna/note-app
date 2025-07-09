document.addEventListener("DOMContentLoaded", () => {
  const form = document.getElementById("note-form");
  const titleInput = document.getElementById("note-title");
  const contentInput = document.getElementById("note-content");
  const notesContainer = document.getElementById("notes-container");

  async function fetchNotes() {
    const res = await fetch("/api/notes", {
      method: "GET",
      credentials: "include", // important if using cookies
    });

    if (!res.ok) return;

    const notes = await res.json();
    notesContainer.innerHTML = "";

    notes.forEach(note => {
      const div = document.createElement("div");
      div.className = "note-card";
      div.innerHTML = `<h3>${note.title}</h3><p>${note.content}</p>`;
      notesContainer.appendChild(div);
    });
  }

  form.addEventListener("submit", async (e) => {
    e.preventDefault();

    const newNote = {
      title: titleInput.value.trim(),
      content: contentInput.value.trim()
    };

    const res = await fetch("/api/notes", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      credentials: "include",
      body: JSON.stringify(newNote)
    });

    if (res.ok) {
      titleInput.value = "";
      contentInput.value = "";
      fetchNotes(); // reload after creating
    }
  });

  fetchNotes();
});
