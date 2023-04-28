let props = {
	msgSocket: undefined,
	fileSocket: undefined,
	fileSocketPort: undefined,
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
};

function msgHandler(data) {
	console.log(data);
}
