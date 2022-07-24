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

  filenameForm.addEventListener("submit", (event) => {
    event.preventDefault();
    apiCallSaveFile(filename.value);
    filenameInput.value = "";
  });

  arduinoConfigForm.addEventListener("submit", (event) => {
    event.preventDefault();
    apiCallSetPort(port.value);
  });
};

document.addEventListener("DOMContentLoaded", init);
