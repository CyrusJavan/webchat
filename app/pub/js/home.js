window.addEventListener("load", function (evt) {
    let output = document.getElementById("output");
    let input = document.getElementById("input");
    let ws;
    let id;
    let token;
    let print = function (message) {
        let d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
    };
    ws = new WebSocket("ws://" + location.hostname + "/chat");
    ws.onopen = function (evt) {
        print("OPEN");
        let toSend = {
            room: "hello",
            action: "join"
        };
        ws.send(JSON.stringify(toSend));
    }
    ws.onclose = function (evt) {
        print("CLOSE");
        ws = null;
    }
    ws.onmessage = function (evt) {
        let d = JSON.parse(evt.data);
        console.log(d)
        if (d.action === "joined") {
            print("JOINED");
            id = d.id;
            token = d.token;
        } else  {
            print("FROM:" + d.id.slice(0, 5) + " MESSAGE:" + d.message);
        }
    }
    ws.onerror = function (evt) {
        print("ERROR: " + evt.data);
    }

    document.getElementById("sendmsg").onclick = sendmsg;
    input.onkeyup = function (e) {
        if (e.key === "Enter") {
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
            id: id,
            token: token
        };
        print("FROM:you" + " MESSAGE:" + toSend.message);
        ws.send(JSON.stringify(toSend));
        input.value = "";
        return false;
    }
});

