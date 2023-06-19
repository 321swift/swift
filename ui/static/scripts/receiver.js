// connect to sender socket
backendSocket = new WebSocket("ws://" + document.location.host + "/receiver");
backendSocket.onmessage = function (evt) {
    handleStatusChange(JSON.parse(evt.data));
};
