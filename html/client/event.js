var selectedchat = "general";
var selectedchatName = "общий";

class Event {
    constructor(type, payload) {
        this.type = type;
        this.payload = payload;
    }
}

class SendMessageEvent {
    constructor(message) {
        this.message = message;
    }
}

class NewMessageEvent {
    constructor(message, sent) {
        this.message = message;
        this.sent = sent;
    }
}

class ChangeChatRoomEvent {
    constructor(name) {
        this.name = name;
    }
}

function routeEvent(event) {

    if (event.type === undefined) {
        alert("no 'type' field in event");
    }
    switch (event.type) {
        case "new_message":
            const messageEvent = Object.assign(new NewMessageEvent, event.payload);
            appendChatMessage(messageEvent);
            break;
        default:
            alert("unsupported message type");
            break;
    }
}

function appendChatMessage(messageEvent) {
    var date = new Date(messageEvent.sent);
  
    const formattedMsg = `${date.toLocaleString()}: ${messageEvent.message}`;
  
    textarea = document.getElementById("chatmessages");
    textarea.innerHTML = textarea.innerHTML + "\n" + formattedMsg;
    textarea.scrollTop = textarea.scrollHeight;
}

function changeChatRoom() {
    var newchat = document.getElementById("chatroom-selection");
    
    if (newchat != null && newchat.value != selectedchat) {
        selectedchat = newchat.value;
        selectedchatName = newchat[newchat.selectedIndex].text;

        textarea = document.getElementById("chatmessages");
        textarea.innerHTML = textarea.innerHTML + "\n" + "Сейчас в чате: " + selectedchatName;
        textarea.scrollTop = textarea.scrollHeight;

        let changeEvent = new ChangeChatRoomEvent(selectedchat);
        sendEvent("change_room", changeEvent);
        textarea = document.getElementById("chatmessages");
        textarea.innerHTML = `Вы поменяли чат на: ${selectedchatName}`;
    }

    return false;
}

function sendMessage() {
    var newmessage = document.getElementById("message");
    if (newmessage != null) {
        let outgoingEvent = new SendMessageEvent(newmessage.value);
        sendEvent("send_message", outgoingEvent)
        newmessage.value = ""
    }
    return false;
}

function sendEvent(eventName, payload) {
    const event = new Event(eventName, payload);
    conn.send(JSON.stringify(event));
}