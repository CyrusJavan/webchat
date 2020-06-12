window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var id;
    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
    };
    ws = new WebSocket("wss://localhost/chat");
    ws.onopen = function(evt) {
        print("OPEN");
        var toSend = {room: "hello", action: "join"};
        ws.send(JSON.stringify(toSend));
    }
    ws.onclose = function(evt) {
        print("CLOSE");
        ws = null;
    }
    ws.onmessage = function(evt) {
        var d = JSON.parse(evt.data);
        console.log("received")
        console.log(d);
        if (d.subscribed) {
            print("JOINED");
            id = d.id;
        } else {
            print("FROM:" + d.id.slice(0,5) + " MESSAGE:" + d.message);
        }
    }
    ws.onerror = function(evt) {
        print("ERROR: " + evt.data);
    }

    document.getElementById("sendmsg").onclick = sendmsg;
    input.onkeyup = function(e){
        if(e.key === "Enter"){
            sendmsg(null);
        }
    }
    function sendmsg(evt) {
        if (!ws || input.value === "") {
            return false;
        }
        let toSend = {
            room: "hello",
            action: "send",
            message: input.value,
            id: id
        };
        print("FROM:you" + " MESSAGE:" + toSend.message);
        ws.send(JSON.stringify(toSend));
        input.value = "";
        return false;
    };
});

