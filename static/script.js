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
    list.insertBefore(li, list.firstChild);
}

function clearRequests() {
    document.getElementById("requestList").innerHTML = '';

    fetch('/api/requests/delete', { method: 'POST' })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to clear requests');
            }
            console.log('All requests cleared');
        })
        .catch(error => console.error('Error:', error));
}

// exported for testability
module.exports = { addRequestToList, fetchInitialData, clearRequests };