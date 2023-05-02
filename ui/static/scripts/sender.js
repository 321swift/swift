// connect to sender socket
backendSocket = new WebSocket("ws://" + document.location.host + "/sender");
backendSocket.onmessage = function (evt) {
	console.log(JSON.parse(evt.data));
};
