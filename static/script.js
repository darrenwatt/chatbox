const wsProtocol = window.location.protocol === "https:" ? "wss" : "ws"; // Determine WebSocket protocol based on HTTP/HTTPS
const input = document.getElementById("messageInput");
const sendButton = document.getElementById("sendButton");
let ws; // WebSocket connection

function connectWebSocket() {
    ws = new WebSocket(`${wsProtocol}://${window.location.host}/ws`);

    ws.onopen = function(event) {
        console.log("WebSocket connection established.");
        const connectButton = document.getElementById("connectButton");
        connectButton.innerText = "Disconnect"; // Change the text here
        connectButton.onclick = disconnectFromChatRoom;
        sendButton.disabled = false;
        input.disabled = false;
    };

    ws.onmessage = function(event) {
        displayMessage(event.data);
    };

    ws.onclose = function(event) {
        console.log("WebSocket connection closed:", event);
        const connectButton = document.getElementById("connectButton");
        connectButton.innerText = "Connect to Chat Room"; // Change the text here
        connectButton.onclick = connectWebSocket;
        sendButton.disabled = true;
        input.disabled = true;
    };
}

function sendMessage() {
    if (input.value.trim() !== "") {
        ws.send(input.value.trim());
        input.value = "";
    }
}

function displayMessage(message) {
    const chatMessages = document.getElementById("chatMessages");
    const messageElement = document.createElement("div");
    messageElement.textContent = message;
    chatMessages.appendChild(messageElement);
}

function disconnectFromChatRoom() {
    console.log("Closing WebSocket connection now.");
    ws.close(1000); // Normal closure
}

document.addEventListener("DOMContentLoaded", () => {
    sendButton.addEventListener("click", sendMessage);
    input.addEventListener("keypress", function(event) {
        if (event.keyCode === 13) {
            event.preventDefault();
            sendMessage();
        }
    });
    document.getElementById("connectButton").addEventListener("click", connectWebSocket, { once: true });
    sendButton.disabled = true;
    input.disabled = true;
});
