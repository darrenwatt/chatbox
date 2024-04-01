
const wsProtocol = window.location.protocol === "https:" ? "wss" : "ws"; // Determine WebSocket protocol based on HTTP/HTTPS

// Get the input field and button
const input = document.getElementById("messageInput");
const sendButton = document.getElementById("sendButton");


function wsHandler() {
    const ws = new WebSocket(`${wsProtocol}://${window.location.host}/ws`);

    ws.onopen = function(event) {
        console.log("WebSocket connection established.");
        // You can add any additional logic here, such as enabling/disabling UI elements
    };

    ws.onmessage = function(event) {
        // Handle incoming messages from the server
        // const message = JSON.parse(event.data);
        const message = event.data;
        displayMessage(message);
    };

    function sendMessage() {
        const input = document.getElementById("messageInput");
        const message = input.value;
        ws.send(message);
        input.value = "";
    }


    function displayMessage(message) {
        const chatMessages = document.getElementById("chatMessages");
        const messageElement = document.createElement("div");
        messageElement.textContent = message; // Set text content directly
        chatMessages.appendChild(messageElement);
    }




    document.getElementById("sendButton").addEventListener("click", sendMessage);

    // Add event listener to input field for "keypress" event
    input.addEventListener("keypress", function(event) {
        // Check if the pressed key is Enter (key code 13)
        if (event.keyCode === 13) {
            // Prevent the default action (form submission)
            event.preventDefault();
            // Call the sendMessage function
            sendMessage();
        }
    });

    sendButton.addEventListener("click", sendMessage);

}

    // Connect to WebSocket server when the "Connect" button is clicked
    document.getElementById("connectButton").addEventListener("click", function() {
        console.log("Attempting to connect to WebSocket server...");
        wsHandler();
    });

