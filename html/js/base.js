/* global jQuery */

//ws://127.0.0.1:8833/ws
var wsServer = 'ws://' +
    window.location.hostname+
    ':' +
    window.location.port+
    '/ws';
var websocket = new WebSocket(wsServer);
websocket.onopen = function (evt) {
    console.log("Connected to WebSocket server.");
    app.echo("Connected to WebSocket server.")
};

websocket.onclose = function (evt) {
    console.log("Disconnected");

    app.echo("ws Disconnected")
};

websocket.onmessage = function (evt) {
    console.log(evt.data);



    app.echo(evt.data)
};

websocket.onerror = function (evt, e) {
    console.log('Error occured: ' + evt.data);
    app.echo(evt.data)
};



app= $('body').terminal(
    function(command) {
        console.log(command)
        websocket.send(command)
    }
);


