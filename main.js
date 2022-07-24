const connect = () => {
  // Connect through websocket
  let ws = new WebSocket("ws://localhost:8080/ws");
  ws.onopen = (e) => {
    const dot = document.getElementById("dot");
    dot.classList.add("connected");
  }
  ws.onclose = (e) => {
    const dot = document.getElementById("dot");
    dot.classList.clear();
    dot.classList.add("disconnected");
  }
  ws.onerror = (e) => {
    // TODO: show red dot and message
  }
  return ws;
};

const apiCallSaveFile = async (filename) => {
  const response = await fetch("/api/save", {
    method: "POST",
    headers: {
      ContentType: "application/json",
    },
    body: JSON.stringify({ filename }),
  });
  // TODO: Check errors here
  const json = await response.json();
};

const apiCallSetPort = async (port) => {
  const response = await fetch("/api/port", {
    method: "PUT",
    headers: {
      ContentType: "application/json",
    },
    body: JSON.stringify({ port }),
  });

  const portErrorMessage = document.querySelector("#arduino-error");
  if (response.status != 200) {
    portErrorMessage.textContent = await response.text();
    return;
  }
  const json = await response.json();
  if (json.ok) {
    portErrorMessage.textContent = "Connected";
  }
};

const init = () => {
  const arduinoConfigForm = document.querySelector("#arduino-form");
  const filename = document.querySelector("#filename-form input[type=text]");

  const filenameForm = document.querySelector("#filename-form");
  const port = document.querySelector("#arduino-form input[type=text]");

  const socket = connect();

  socket.onmessage = (e) => {
    const result = JSON.parse(e.data);
    switch (result.cmd) {
      case "setPort": {
        const arduinoError = document.querySelector("#arduino-error");
        arduinoError.textContent = result.ok ? "Connected" : result.error;
        break;
      }
      case "saveFile": {
        console.log("saveFile", result);
        break;
      }
    }
  };

  filenameForm.addEventListener("submit", (event) => {
    event.preventDefault();
    socket.send(JSON.stringify({ cmd: "saveFile", args: [filename.value] }));
    filename.value = "";
  });

  arduinoConfigForm.addEventListener("submit", (event) => {
    event.preventDefault();
    socket.send(JSON.stringify({ cmd: "setPort", args: [port.value] }));
  });
};

document.addEventListener("DOMContentLoaded", init);
