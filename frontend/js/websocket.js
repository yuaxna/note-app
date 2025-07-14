let ws = null;

function connectWebSocket() {
  if (ws && ws.readyState === WebSocket.OPEN) return;

  ws = new WebSocket("ws://localhost:8080/api/ws");

  ws.onopen = () => console.log("✅ WebSocket connected");

  ws.onmessage = (e) => {
    const msg = JSON.parse(e.data);
    switch (msg.action) {
      case "create":
      case "edit":
      case "delete":
        handleLiveUpdate(msg);
        break;
      default:
        console.warn("Unknown action:", msg.action);
    }
  };

  ws.onclose = () => {
    console.warn("🔌 WebSocket disconnected. Retrying in 3s...");
    setTimeout(connectWebSocket, 3000);
  };

  ws.onerror = (err) => console.error("WebSocket error:", err);
}

function handleLiveUpdate(msg) {
  console.log("📬 Live update received:", msg);

  // Simply re-fetch all notes when anything changes
  if (typeof window.fetchNotes === "function") {
    window.fetchNotes(); // <-- this refreshes the list
  } else {
    console.warn("⚠️ fetchNotes is not defined yet.");
  }
}

function sendNoteUpdate(action, noteId, title, content, sender) {
  if (ws && ws.readyState === WebSocket.OPEN) {
    const payload = {
      action,
      note_id: noteId,
      title,
      content,
      sender,
    };
    ws.send(JSON.stringify(payload));
  } else {
    console.warn("⚠️ Cannot send WebSocket message — not connected");
  }
}

// Auto-connect on page load
document.addEventListener("DOMContentLoaded", connectWebSocket);

// Make available globally (if needed elsewhere)
window.connectWebSocket = connectWebSocket;
window.sendNoteUpdate = sendNoteUpdate;
