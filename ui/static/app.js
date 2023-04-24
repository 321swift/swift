let props = {
	socket: undefined,
	fileSocket: undefined,
	fileSocketPort: undefined,
	statusQueue: [""],
};

let sendBtn = document.getElementById("sendBtn");
let receiveBtn = document.getElementById("receiveBtn");
let statusArea = document.getElementById("status");
let fileInput = document.getElementById("fileInput");
console.log(fileInput);

let updateEvent = new Event("update", {});

fileInput.addEventListener("change", (event) => {
	const reader = new FileReader();
	reader.addEventListener("load", (event) => {
		const result = event.target.result;
		// Do something with result
		console.log(result);
		if (props.socket != undefined) {
			props.fileSocket.send(result);
		}
	});

	reader.addEventListener("progress", (event) => {
		if (event.loaded && event.total) {
			const percent = (event.loaded / event.total) * 100;
			console.log(`Progress: ${Math.round(percent)}`);
		}
	});

	reader.readAsArrayBuffer(event.target.files[0]);
});

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
			// parse websocket message as JSON
			const eventData = JSON.parse(evt.data);
			if (JSON.parse(eventData).hasOwnProperty("fileSocket")) {
				props.fileSocketPort = JSON.parse(eventData).fileSocket;
				if (props.fileSocket == undefined) {
					setupFileSocket();
				}
			}
			update(eventData);
		};
	} else {
		alert("Not supporting websockets, \nPlease use a different browser");
	}
}

function setupFileSocket() {
	let conn = new WebSocket("ws://" + document.location.host + "/ws");

	// Onopen
	conn.onopen = function (evt) {
		update("connected to fileSocket");
		props.fileSocket = conn;
	};
	conn.onclose = function (evt) {
		// Set disconnected
		update("file socket closed");
		props.fileSocket = undefined;
	};
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
