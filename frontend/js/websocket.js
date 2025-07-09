const ws = new WebSocket("ws://localhost:8080/ws");

ws.onopen = () => {
  console.log("WebSocket connected");
};

ws.onmessage = (event) => {
  console.log("Message from server:", event.data);
  // Here, update your notes UI live based on received message
};

ws.onclose = () => {
  console.log("WebSocket disconnected");
};

ws.onerror = (error) => {
  console.error("WebSocket error:", error);
};

function sendNoteUpdate(note) {
  if (ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({
      ...note,
      sender: currentUsername // you'd fetch this from session
    }));
  }
}
