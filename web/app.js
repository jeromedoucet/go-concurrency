"use strict";

var serversocket = new WebSocket("ws://localhost:4444/websocket");

// Write message on receive
serversocket.onmessage = function(e) {
    var notification = JSON.parse(e.data);
    var id = document.createElement("td");
    id.appendChild(document.createTextNode(notification.PlayerId));

    var state = document.createElement("td");
    state.appendChild(document.createTextNode(notification.Type === 'registrate' ? 'ok' : 'ko'));

    var Rate = document.createElement("td");
    Rate.appendChild(document.createTextNode(notification.Rate));

    var Score = document.createElement("td");
    Score.appendChild(document.createTextNode(notification.Score));

    var tr = document.createElement("tr");
    tr.id = e.PlayerId
    tr.appendChild(id);
    tr.appendChild(state);
    tr.appendChild(Rate);
    tr.appendChild(Score);
    document.getElementById("table-body").appendChild(tr);
};
