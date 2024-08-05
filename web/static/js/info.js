//  Variables
var SUB_INFO_TOPIC = "info";
var SUB_INFO_AREA = "txt_sysInfo";
var SUB_EVENTS_TOPIC = "status/io/sensorhub/events/#";

// called when a message arrives
const topicHandlers_Info = {
    info: handleInfoMessage,
    status: {
        io: {
            sensorhub: {
                events: {
                    crash: handleCrashEventMessage,
                    PSE: handlePSEEventMessage,
                    RTC: handleRTCEventMessage,
                    Thermal: handleThermalEventMessage,
                    USB: handleUSBEventMessage
                }
            }
        }
    }
};

function onMessageArrived_Info(message) {
    if (document.readyState !== 'complete') {
        console.log('Document is not loaded yet, onMessageArrived_Info()');
        return;
    }
    console.log('onMessageArrived_Info() called!');

    const fullTopic = message.destinationName;
    const payload = message.payloadString;
    const topics = fullTopic.split("/");

    let handler = topicHandlers_Info;
    for (let i = 0; i < topics.length; i++) {
        handler = handler[topics[i]];
        if (handler === undefined) {
            console.log(`No topic matched! Topic = ${topics.slice(0, i + 1).join('/')}`);
            return;
        }
        if (typeof handler === 'function') {
            handler(payload);
            return;
        }
    }
    console.log(`onMessageArrived_Info:No topic matched! Topic = ${fullTopic}`);
}

function handleCrashEventMessage(payloada) {
    const formattedPayload = payloada.replace(/\n/g, '<br>');
    var label = document.querySelector("label[for='crashStatus']");
    //check if the label exists
    if (label == null) {
        return;
    }
    label.textContent = formattedPayload;
    // Parse payload to find "int=0" or "int=1"
    var intMatch = payloada.match(/int=(\d)/);
    if (intMatch) {
        var intValue = intMatch[1];
        if (intValue === '0') {
            crashIcon.style.color = 'gray'; // Change the color to gray
        } else if (intValue === '1') {
            crashIcon.style.color = '#ff0000'; // Change the color to red
        }
    } else {
        // no match, set to gray
        crashIcon.style.color = 'gray'; // Change the color to gray
    }
}

function handlePSEEventMessage(payloada) {
    const formattedPayload = payloada.replace(/\n/g, '<br>');
    var label = document.querySelector("label[for='pseStatus']");
    //check if the label exists
    if (label == null) {
        return;
    }
    if (payloada == 'on') {
        label.textContent = "PSE: ON";
        // Change the color to Green
        pseIcon.style.color = '#00ff00';
    } else if (payloada == 'off') {
        label.textContent = "PSE: OFF";
        // Change the color to white
        pseIcon.style.color = '#c2c5db';
    }
}

let timerId = null;  // Declare timerId
let startTime = new Date();  // Get the current system time
window.startTime = startTime; // Make startTime global

function handleRTCEventMessage(payloada) {
    // Stop the exist timer for tick function
    if (timerId) {
        clearInterval(timerId);
    }
    // Parse the payload as time and set it as the start time
    startTime = new Date(payloada);
    const formattedPayload = payloada.replace(/\n/g, '<br>');
    var label = document.querySelector("label[for='rtcStatus']");
    //check if the label exists
    if (label == null) {
        return;
    }
    label.textContent = startTime.toLocaleString();
    /* if you need to count by browser itself, enable this code 
       Call the tick function every 1 second (1000 milliseconds)*/
    // timerId = setInterval(tick, 1000);
}

function tick() {
    startTime.setSeconds(startTime.getSeconds() + 1);
    console.log('Current time: ' + startTime.toString());
    // let rtcStatus label show system time
    var label = document.querySelector("label[for='rtcStatus']");
    //check if the label exists
    if (label == null) {
        return;
    }
    label.textContent = startTime.toLocaleString();
}

function handleThermalEventMessage(payloada) {
    const formattedPayload = payloada.replace(/\n/g, '<br>');
    var label = document.querySelector("label[for='thermalStatus']");
    //check if the label exists
    if (label == null) {
        return;
    }
    label.textContent = formattedPayload;

    /*************************************************************************
    // if payloada value > 80, change the color to red, else change to green
    // depends on the payloada value, the following code may need to be changed
    if (parseInt(payloada) > 80) {
        thermalIcon.style.color = '#ff0000'; // Change the color to red
    } else {
        thermalIcon.style.color = '#00ff00'; // Change the color to green
    }
    ****************************************************************************/
}

function handleUSBEventMessage(payloada) {
    console.log('handleUSBEventMessage() called!');
    const formattedPayload = payloada.replace(/\n/g, '<br>');
    var label = document.querySelector("label[for='usbPlug']");
    //check if the label exists
    if (label == null) {
        return;
    }
    if (payloada == 'in') {
        label.textContent = "USB Connected";
        // Change the color to Blue        
        usbIcon.style.color = '#0000ff';
    } else if (payloada == 'out') {
        label.textContent = "USB Disconnected";
        usbIcon.style.color = '#c2c5db'; // Change the color to white
    }

}

function handleInfoMessage(payloada) {
    console.log('handleInfoMessage() called!');
    const formattedPayload = payloada.replace(/\n/g, '<br>');
    const cardBody = document.getElementById(SUB_INFO_AREA);
    // check if the cardBody exists
    if (cardBody == null) {
        console.log("cardBody not found");
        return;
    }
    cardBody.innerHTML = formattedPayload;
    cardBody.scrollTop = cardBody.scrollHeight;
}
