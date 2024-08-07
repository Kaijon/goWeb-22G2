//  Variables
//var SUB_FOTA_TOPIC = "fota/#";
var SUB_FOTA_TOPIC = "status/fota/#";
var SUB_FOTA_AREA = "progress_info";
//var SUB_EVENTS_TOPIC = "status/io/sensorhub/events/#";

// called when a message arrives
// 1.uboot 2.env 3.Image 4.dtb 4.rootfs 5.flash
const topicHandlers_Fota = {
    fota: handleFotaMessage,
    status: {
        fota: {
            1: handleFota1Message,
            2: handleFota2Message,
            3: handleFota1Message,
            4: handleFota1Message,
            5: handleFota1Message
        }
    }
};

function onMessageArrived_Fota(message) {
    console.log(`onMessageArrived_Fota`)
    if (document.readyState !== 'complete') {
        console.log('Document is not loaded yet, onMessageArrived_Fota()');
        return;
    }
    console.log('onMessageArrived_Fota() called!');

    const fullTopic = message.destinationName;
    const payload = message.payloadString;
    const topics = fullTopic.split("/");

    let handler = topicHandlers_Fota;
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

//function handleFotaMessage(payloada) {
//    console.log('handleInfoMessage() called!');
//    const formattedPayload = payloada.replace(/\n/g, '<br>');
//    const cardBody = document.getElementById(SUB_FOTA_AREA);
//    // check if the cardBody exists
//    if (cardBody == null) {
//        console.log("cardBody not found");
//        return;
//    }
//    cardBody.innerHTML = formattedPayload;
//    cardBody.scrollTop = cardBody.scrollHeight;
//}


function handleFotaMessage(payload) {
    console.log('handleInfoMessage() called!');
    // check if the cardBody exists
    const cardBody = document.getElementById(SUB_FOTA_AREA);
    if (cardBody == null) {
        console.log("cardBody not found");
        return;
    }
    let formattedPayload = '';
    let parsedPayload;
    let fotaMessages = [];

    try {
        parsedPayload = JSON.parse(payload);
    } catch (error) {
        console.log(`Fail to parse JSON text`);
        cardBody.innerHTML = "Fail to parse JSON text";
        cardBody.scrollTop = cardBody.scrollHeight;
        return; 
    }

    fotaMessages.push(parsedPayload);
    fotaMessages.forEach((message, index) => {
        formattedPayload += `<strong>Message ${index + 1}:</strong><br>`;
        
        // Create a map to store key-value pairs
        const messageMap = new Map(Object.entries(message));

        messageMap.forEach((value, key) => {
            switch (key) {
                case 'stage':
                    formattedPayload += `Stage: ${value}<br>`;
                    break;
                case 'percentage':
                    formattedPayload += `Percentage: ${value}%<br>`;
                    break;
                case 'status':
                    formattedPayload += `Status: ${value}<br>`;
                    break;
                default:
                    formattedPayload += `${key}: ${value}<br>`;
                    break;
            }
        });

        formattedPayload += '<br>';
    });


    cardBody.innerHTML = formattedPayload;
    cardBody.scrollTop = cardBody.scrollHeight;
}


function handleFota1Message(payload) {
    console.log('handleFota1Message() called!');
    // check if the cardBody exists
    const cardBody = document.getElementById(SUB_FOTA_AREA);
    if (cardBody == null) {
        console.log("cardBody not found");
        return;
    }
    let formattedPayload = '';
    let parsedPayload;
    let fotaMessages = [];

    try {
        parsedPayload = JSON.parse(payload);
    } catch (error) {
        console.log(`Fail to parse JSON text`);
        cardBody.innerHTML = "Fail to parse JSON text";
        cardBody.scrollTop = cardBody.scrollHeight;
        return; 
    }

    fotaMessages.push(parsedPayload);
    fotaMessages.forEach((message, index) => {
        formattedPayload += `<strong>Message ${index + 1}:</strong><br>`;
        
        // Create a map to store key-value pairs
        const messageMap = new Map(Object.entries(message));

        messageMap.forEach((value, key) => {
            switch (key) {
                case 'percentage':
                    formattedPayload += `Percentage: ${value}%<br>`;
                    break;
                case 'status':
                    formattedPayload += `Status: ${value}<br>`;
                    break;
                default:
                    formattedPayload += `${key}: ${value}<br>`;
                    break;
            }
        });

        formattedPayload += '<br>';
    });


    cardBody.innerHTML = formattedPayload;
    cardBody.scrollTop = cardBody.scrollHeight;
}


function handleFota2Message(payload) {
    console.log('handleFota2Message() called!');
    // check if the cardBody exists
    const cardBody = document.getElementById(SUB_FOTA_AREA);
    if (cardBody == null) {
        console.log("cardBody not found");
        return;
    }
    let formattedPayload = '';
    let parsedPayload;
    let fotaMessages = [];

    try {
        parsedPayload = JSON.parse(payload);
    } catch (error) {
        console.log(`Fail to parse JSON text`);
        cardBody.innerHTML = "Fail to parse JSON text";
        cardBody.scrollTop = cardBody.scrollHeight;
        return; 
    }

    fotaMessages.push(parsedPayload);
    fotaMessages.forEach((message, index) => {
        formattedPayload += `<strong>Message ${index + 1}:</strong><br>`;
        
        // Create a map to store key-value pairs
        const messageMap = new Map(Object.entries(message));

        messageMap.forEach((value, key) => {
            switch (key) {
                case 'percentage':
                    formattedPayload += `Percentage: ${value}%<br>`;
                    break;
                case 'status':
                    formattedPayload += `Status: ${value}<br>`;
                    break;
                default:
                    formattedPayload += `${key}: ${value}<br>`;
                    break;
            }
        });

        formattedPayload += '<br>';
    });


    cardBody.innerHTML = formattedPayload;
    cardBody.scrollTop = cardBody.scrollHeight;
}



























