function getAPIUrl(endpoint) {
    const idInstance = document.getElementById('idInstance').value;
    const apiTokenInstance = document.getElementById('apiTokenInstance').value;
    return `https://api.green-api.com/waInstance${idInstance}/${endpoint}/${apiTokenInstance}`;
}

async function fetchAPI(endpoint, options = {}) {
    try {
        const response = await fetch(getAPIUrl(endpoint), {
            ...options,
            headers: {
                'Content-Type': 'application/json'
            }
        });
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        return await response.json();
    } catch (error) {
        console.error('Error fetching API:', error);
        appendResponse('Error', error.message);
        throw error;
    }
}

function appendResponse(method, result) {
    const apiResponses = document.getElementById('apiResponses');
    const currentContent = apiResponses.value;
    const newContent = `${method}:\n${JSON.stringify(result, null, 2)}\n\n`;
    apiResponses.value = currentContent + newContent;
}

async function getSettings() {
    try {
        const result = await fetchAPI('getSettings');
        appendResponse('getSettings', result);
    } catch (error) {
        appendResponse('getSettings Error', error.message);
    }
}

async function getStateInstance() {
    try {
        const result = await fetchAPI('getStateInstance');
        appendResponse('getStateInstance', result);
    } catch (error) {
        appendResponse('getStateInstance Error', error.message);
    }
}

async function sendMessage() {
    const phoneNumber = document.getElementById('phoneNumber').value;
    const messageText = document.getElementById('messageText').value;

    try {
        const result = await fetchAPI('sendMessage', {
            method: 'POST',
            body: JSON.stringify({
                phoneNumber,
                messageText
            })
        });
        appendResponse('sendMessage', result);
    } catch (error) {
        appendResponse('sendMessage Error', error.message);
    }
}

async function sendFileByUrl() {
    const phoneNumber = document.getElementById('filePhoneNumber').value;
    const fileUrl = document.getElementById('fileUrl').value;

    try {
        const result = await fetchAPI('sendFileByUrl', {
            method: 'POST',
            body: JSON.stringify({
                phoneNumber,
                fileUrl
            })
        });
        appendResponse('sendFileByUrl', result);
    } catch (error) {
        appendResponse('sendFileByUrl Error', error.message);
    }
}
