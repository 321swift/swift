let props = {
	socket: undefined,
	statusQueue: [""],
};

let sendBtn = document.getElementById("sendBtn");
let receiveBtn = document.getElementById("receiveBtn");
let statusArea = document.getElementById("status");

let updateEvent = new Event("update", {});

receiveBtn.onclick = (e) => {
	if (props.socket != undefined) {
		console.log("sending json");
		props.socket.send(
			JSON.stringify({
				role: "client",
			})
		);
	}
};
sendBtn.onclick = (e) => {
	if (props.socket != undefined) {
		console.log("sending json");
		props.socket.send(
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
			update("connected to backend");
			props.socket = conn;
		};

		conn.onclose = function (evt) {
			// Set disconnected
			update("socket closed");
			props.socket = undefined;
		};

		// Add a listener to the onmessage event
		conn.onmessage = function (evt) {
			console.log(evt);
			// parse websocket message as JSON
			const eventData = JSON.parse(evt.data);

			update(eventData);
		};
	} else {
		alert("Not supporting websockets, \nPlease use a different browser");
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
		newChild.textContent = props.statusQueue[props.statusQueue.length - 1];
		statusArea.appendChild(newChild);
	});
}

function update(data) {
	if (data != undefined) {
		props.statusQueue.push(data);
		statusArea.dispatchEvent(updateEvent);
	}
}
