class Queue {
	constructor() {
		this.queue = [];
	}

	dequeue() {
		let a = this.queue.shift();
		return a;
	}
	enqueue(item) {
		this.queue.push(item);
	}

	isEmpty() {
		return this.queue.length == 0;
	}

	get() {
		return this.queue;
	}
}

let enqueueEvent = new Event("enqueue");
let props = {
	msgSocket: undefined,
	fileSocket: undefined,
	fileSocketPort: undefined,
	connectedMessage: `Connection made successfully <br\> please select the files you wish to send `,
	fileInputArea: null,
	fileInput: null,
	statusArea: null,
	fileQueue: new Queue(),
};

function connectMsgSocket() {
	if (!window["WebSocket"]) {
		alert(
			"This application uses web sockets in its operations.\
      \nYour browser unfortunately does not support web sockets. \
      \nPlease try using a different browser"
		);
		return;
	}

	let conn = new WebSocket("ws://" + document.location.host + "/ws");

	conn.onopen = function (evt) {
		props.msgSocket = conn;
		props.msgSocket.send(JSON.stringify({ role: "server" }));
	};

	conn.onclose = function (evt) {
		props.msgSocket = undefined;
	};

	conn.onmessage = function (evt) {
		msgHandler(JSON.parse(evt.data));
	};
}

window.onload = function () {
	connectMsgSocket();
	props.fileInputArea = document.getElementById("fileInputArea");
	props.statusArea = document.getElementById("statusArea");

	window.addEventListener("enqueue", (event) => {
		while (!props.fileQueue.isEmpty()) {
			if (props.fileSocket != undefined) {
				props.fileSocket.send(props.fileQueue.dequeue());
			}
		}
	});
};

function updateUI(data) {
	switch (data.status) {
		case "ok":
			break;
		case "assuming server role":
			break;

		default:
			// console.log(data);
			// console.log(data.status);
			if (data.status.startsWith("Connection made:")) {
				statusArea.innerHTML = props.connectedMessage;
				if (props.fileSocketPort != undefined) {
					setupFileSocket();
				}
			}
	}
}

function msgHandler(data) {
	if (!isJsonString(data)) {
		console.log("this is data \n", data);
		return;
	}

	switch (Object.keys(JSON.parse(data))[0]) {
		case "status":
			updateUI(JSON.parse(data));
			break;
		case "log":
			// console.log(data);
			break;
		case "fileSocket":
			console.log(data);
			props.fileSocketPort = JSON.parse(data).fileSocket;
			break;
		case "file":
			console.log("this is data   \n", data);
		// console.log(JSON.parse(JSON.parse(data).file).Data);
	}
}

function setupFileSocket() {
	let conn = new WebSocket(
		"ws://" + document.location.hostname + ":" + props.fileSocketPort + "/file"
	);

	// Onopen
	conn.onopen = function (evt) {
		console.log("connected to fileSocket");
		props.fileSocket = conn;

		setupFileInput();
	};
	conn.onclose = function (evt) {
		// Set disconnected
		console.log("file socket closed");
		props.statusArea.innerHTML = "Backend disconnected <br /> please restart application";
		props.fileSocket = undefined;
		props.fileInput.remove();
	};
}

function isJsonString(str) {
	try {
		JSON.parse(str);
	} catch (e) {
		return false;
	}
	return true;
}

function readFile(file) {
	return new Promise((resolve, reject) => {
		let fr = new FileReader();
		fr.onload = (event) => {
			// props.fileQueue.enqueue({
			// 	Filename: file.name,
			// 	Data: event.target.result,
			// });
			props.fileQueue.enqueue(event.target.result);
			resolve();
		};
		fr.readAsArrayBuffer(file);
	});
}

function setupFileInput() {
	let fileInput = document.createElement("input");
	fileInput.setAttribute("type", "file");
	fileInput.setAttribute("multiple", "");
	props.fileInput = fileInput;
	props.fileInputArea.appendChild(fileInput);

	fileInput.addEventListener("change", (event) => {
		let filePromises = [];
		Array.from(event.target.files).forEach((file) => {
			filePromises.push(readFile(file));
		});
		Promise.all(filePromises).then(() => {
			// console.log(props.fileQueue.get());
			window.dispatchEvent(enqueueEvent);
		});
	});
}
