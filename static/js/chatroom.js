var g = {};

function pushChatLog(jsonMessage) {
  console.log("PUSH MESSAGE", jsonMessage);
  let needScroll = false;
  if (g.chatLog.scrollHeight - g.chatLog.scrollTop === g.chatLog.clientHeight) {
    needScroll = true;
  }
  // if (jsonMessage.sender === g.displayName) {
  //   needScroll = true;
  // }

  var item = document.createElement("div");

  if (jsonMessage.sender === "SYSTEM") {
    item.innerHTML =
      '<b><span style="color: red;">SYSTEM: </span>' +
      jsonMessage.message +
      "</b>";
  } else {
    item.innerHTML =
      "<span><b>" + jsonMessage.sender + ":</b> </span>" + jsonMessage.message;
  }

  if (jsonMessage.sender) jsonMessage.message;
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

  let jsonStr = null;
  if (g.message.value.startsWith("!")) {
    jsonStr = JSON.stringify({
      type: "_SYSCOMMAND",
      message: g.message.value,
    });
  } else {
    jsonStr = JSON.stringify({
      type: "MESSAGE",
      message: g.message.value,
    });
  }

  console.log("Sending ", jsonStr);
  g.conn.send(jsonStr);
  g.message.value = "";

  return false;
}

function makeWebSocket() {
  g.conn = new WebSocket("wss://" + document.location.host + "/ws");

  g.conn.onclose = function (evt) {
    pushChatLog({ sender: "SYSTEM", message: "Connection closed." });
    g.message.disabled = true;
    document.getElementById("submitBtn").disabled = true;
  };
  g.conn.onmessage = function (evt) {
    var messages = evt.data.split("\n");
    for (var i = 0; i < messages.length; i++) {
      console.log("Received message ", JSON.parse(messages[i]));
      handleIncomingMessage(messages[i]);
    }
  };
  g.conn.onopen = (evt) => {
    console.log("Connected to server.");

    console.log(
      'Sending { type: "_SYSCOMMAND", message: "!get_display_name" }'
    );
    g.conn.send(
      JSON.stringify({
        type: "_SYSCOMMAND",
        message: "!get_display_name",
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
    switch (json.type) {
      case "_SYSCOMMAND":
        handleSysCommand(json);
        break;
      case "MESSAGE":
        pushChatLog(json);
        break;
      case "_SYSMESSAGE":
        pushChatLog({sender: "SYSTEM", message: json.message});
        break;
      default:
        console.log("Incoming message type mismatch");
    }
  });
};

const handleSysCommand = (jsonObj) => {
  switch (jsonObj.message) {
    case "!get_display_name":
      g.displayName = jsonObj.response;
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
