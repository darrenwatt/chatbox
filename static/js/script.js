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
        console.log('Message from server ', event.data);
        if (event.data.startsWith('status:')) {  // Check if the message is a room status update
            document.getElementById('room-status').innerText = event.data.substring(7);  // Remove "status:" prefix
        } else if (event.data.startsWith('chat:')) {     
        displayMessage(event.data.substring(5));
        }
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
    //chatMessages.appendChild(messageElement);
    if (chatMessages.firstChild) {
        chatMessages.insertBefore(messageElement, chatMessages.firstChild);
    } else {
        chatMessages.appendChild(messageElement);
    }

}

function disconnectFromChatRoom() {
    document.getElementById('room-status').innerText = ""
    console.log("Closing WebSocket connection now.");
    ws.close(1000); // Normal closure
}

document.addEventListener("DOMContentLoaded", () => {

    // const toggleLink = document.querySelector('.toggle-link');
    // const safetyText = document.querySelector('.safety-text');

    // toggleLink.addEventListener('click', function (event) {
    //     event.preventDefault(); // Prevent the default link behavior
    //     safetyText.classList.toggle('open'); // Toggle the 'open' class
    //     toggleLink.textContent = safetyText.classList.contains('open') ? 'Close' : 'Safety Reminder';
    //   });
   

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
