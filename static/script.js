function fetchInitialData() {
    fetch('/api/requests/all')
        .then(response => response.json())
        .then(data => {
            data.forEach(request => {
                addRequestToList(request);
            });
        })
        .catch(error => console.error('Error fetching initial data:', error));
}

function addRequestToList(message) {
    var li = document.createElement("li");
    li.className = 'request-detail';
    li.innerHTML = message.split('\n').map(function(detail) {
        return '<span>' + detail.split(': ')[0] + ':</span> ' + detail.split(': ').slice(1).join(': ');
    }).join('<br>');
    var list = document.getElementById("requestList");
    // Insert the new element at the beginning of the list
    list.insertBefore(li, list.firstChild);
}

document.getElementById("clearRequests").addEventListener("click", function() {
    // Clear the list in the frontend
    document.getElementById("requestList").innerHTML = '';

    // Send a request to the backend to clear BoltDB
    fetch('/api/requests/delete', { method: 'POST' })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to clear requests');
            }
            console.log('All requests cleared');
        })
        .catch(error => console.error('Error:', error));
});

var ws = new WebSocket("ws://localhost:8000/ws");

ws.onmessage = function(event) {
    addRequestToList(event.data);
};

ws.onopen = function(event) {
    ws.send("Connected");
};

document.addEventListener('DOMContentLoaded', fetchInitialData);