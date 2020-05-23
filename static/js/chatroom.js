var g = {};

function pushChatLog(jsonMessage) {
  console.log("PUSH MESSAGE", jsonMessage);
  let needScroll = false;
  if (g.chatLog.scrollHeight - g.chatLog.scrollTop === g.chatLog.clientHeight) {
    needScroll = true;
  }
  if (jsonMessage.Sender === g.displayName) {
    needScroll = true;
  }

  var item = document.createElement("div");

  if (jsonMessage.Sender === "SYSTEM") {
    item.innerHTML =
      '<b><span style="color: red;">SYSTEM: </span>' +
      jsonMessage.Message +
      "</b>";
  } else {
    item.innerHTML =
      "<span><b>" + jsonMessage.Sender + ":</b> </span>" + jsonMessage.Message;
  }

  if (jsonMessage.Sender) jsonMessage.message;
  g.chatLog.appendChild(item);

  if (needScroll) {
    g.chatLog.scrollTop = g.chatLog.scrollHeight;
  }
}

function sendMessage() {
  if (!g.conn) {
    return false;
  }
  if (!g.message.value) {
    return false;
  }

  const jsonStr = JSON.stringify({
    Type: "MESSAGE",
    Message: g.message.value,
    Sender: g.displayName,
  });
  console.log("Sending ", jsonStr);
  g.conn.send(jsonStr);
  g.message.value = "";

  return false;
}

function makeWebSocket() {
  g.conn = new WebSocket("wss://" + document.location.host + "/ws");

  g.conn.onclose = function (evt) {
    pushChatLog({ Sender: "SYSTEM", Message: "Connection closed." });
    g.message.disabled = true;
    document.getElementById("submitBtn").disabled = true;
  };
  g.conn.onmessage = function (evt) {
    var messages = evt.data.split("\n");
    for (var i = 0; i < messages.length; i++) {
      console.log("Received message " + messages[i], JSON.parse(messages[i]));
      handleIncomingMessage(messages[i]);
    }
  };
  g.conn.onopen = (evt) => {
    console.log("Connected to server.");

    console.log(
      'Sending { Type: "_SYSCOMMAND", Message: "!get_display_name" }'
    );
    g.conn.send(
      JSON.stringify({
        Type: "_SYSCOMMAND",
        Message: "!get_display_name",
        Sender: g.displayName,
      })
    );
  };
}

const handleIncomingMessage = (message) => {
  let jsonObj = JSON.parse(message);
  if (!Array.isArray(jsonObj)) {
    jsonObj = [jsonObj];
  }

  console.log(jsonObj);

  jsonObj.forEach((json) => {
    switch (json.Type) {
      case "_SYSCOMMAND":
        handleSysCommand(json);
        break;
      case "MESSAGE":
        pushChatLog(json);
        break;
      default:
        console.log("Incoming message type mismatch");
    }
  });
};

const handleSysCommand = (jsonObj) => {
  switch (jsonObj.Message) {
    case "!get_display_name":
      g.displayName = jsonObj.Response;
      console.log("Display name set to " + g.displayName);
      break;
    default:
      console.log("Incoming SYSCOMMAND type mismatch");
  }
};

function init() {
  g.message = document.getElementById("message");
  g.chatLog = document.getElementById("chatlog");

  document.getElementById("form").onsubmit = sendMessage;

  if (window["WebSocket"]) {
    makeWebSocket();
  } else {
    pushChatLog("<b>Your browser does not support WebSockets.</b>");
  }
}

window.onload = init;
