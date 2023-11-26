function addRequestToList(message) {
    var list = document.getElementById("requestList");
    list.innerHTML = message + list.innerHTML
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
module.exports = { addRequestToList, clearRequests };