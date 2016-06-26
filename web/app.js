"use strict";

var serversocket = new WebSocket("ws://localhost:4444/websocket");

// Write message on receive
serversocket.onmessage = function(e) {
    var notification = JSON.parse(e.data);
    console.log(notification);
    if (notification.Type === 'registrate' || notification.Type === 'unregistrate') {
        var id = document.createElement("td");
        id.appendChild(document.createTextNode(notification.PlayerId));
        id.id = notification.PlayerId + '_playerId';

        var state = document.createElement("td");
        state.appendChild(document.createTextNode(notification.Type === 'registrate' ? 'ok' : 'ko'));
        state.id = notification.PlayerId + '_state';

        var rate = document.createElement("td");
        rate.appendChild(document.createTextNode(notification.Rate));
        rate.id = notification.PlayerId + '_rate';

        var score = document.createElement("td");
        score.appendChild(document.createTextNode(notification.Score));
        score.id = notification.PlayerId + '_score';

        var tr = document.createElement("tr");
        tr.id = notification.PlayerId
        tr.appendChild(id);
        tr.appendChild(state);
        tr.appendChild(rate);
        tr.appendChild(score);
        document.getElementById("table-body").appendChild(tr);
    } else {
        document.getElementById(notification.PlayerId + '_score').innerHTML = notification.Score + '$';
        document.getElementById(notification.PlayerId + '_rate').innerHTML = notification.Rate + '$ / sec';
    }
};
