let socket;
let statusQueue = [""];

let sendBtn = document.getElementById("sendBtn");
let receiveBtn = document.getElementById("receiveBtn");
let statusArea = document.getElementById("status");

let updateEvent = new Event("update", {});

sendBtn.onclick = (e) => {
	if (socket != undefined) {
		console.log("sending json");
		socket.send(
			JSON.stringify({
				role: "server",
			})
		);
	}
};

function connectWebsocket() {
	// Check if the browser supports WebSocket
	if (window["WebSocket"]) {
		console.log("supports websockets");
		// Connect to websocket using OTP as a GET parameter
		let conn = new WebSocket("ws://" + document.location.host + "/ws");
		// Onopen
		conn.onopen = function (evt) {
			statusQueue.push("socket open");
			statusArea.dispatchEvent(updateEvent);
			socket = conn;
		};

		conn.onclose = function (evt) {
			// Set disconnected
			statusQueue.push("socket closed");
			statusArea.dispatchEvent(updateEvent);
		};

		// Add a listener to the onmessage event
		conn.onmessage = function (evt) {
			console.log(evt);
			// parse websocket message as JSON
			const eventData = JSON.parse(evt.data);

			// add event data to queue
			statusQueue.push(eventData);
			statusArea.dispatchEvent(updateEvent);
			// // Let router manage message
			// routeEvent(event);
		};
	} else {
		alert("Not supporting websockets");
	}
}
/**
 * Once the website loads
 * */
window.onload = function () {
	// Apply our listener functions to the submit event on both forms
	// we do it this way to avoid redirects
	connectWebsocket();
	handleUpdates();
};

function handleUpdates() {
	statusArea.addEventListener("update", (event) => {
		let newChild = document.createElement("p");
		newChild.textContent = statusQueue[statusQueue.length - 1];
		statusArea.appendChild(newChild);
	});
}
