let props = {
	msgSocket: undefined,
	fileSocket: undefined,
	fileSocketPort: undefined,
};
// connect to the websocket
function connectMsgSocket() {
	if (!window["Websocket"]) {
		alert(
			"This application uses web sockets in its operations.\
      \nYour browser unfortunately does not support web sockets. \
      \nPlease try using a different browser"
		);
		return;
	}

	let conn = new Websocket("ws://" + document.location.host + "/ws");

	conn.onopen = function (evt) {
		props.msgSocket = conn;
	};

	conn.onclose = function (evt) {
		props.msgSocket = undefined;
	};

	conn.onmessage = function (evt) {
		msgHandler(JSON.parse(evt.data));
	};
}

window.onload = function () {
	connectWebsocket();
};

function msgHandler(data) {
	console.log(data);
}
