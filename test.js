const { addRequestToList, fetchInitialData, clearRequests } = require('./static/script');

describe('DOM manipulation tests', () => {
    beforeEach(() => {
        document.body.innerHTML = '<ul id="requestList"></ul>';
        require('./static/script');
    });

    test('addRequestToList adds a new item to the list', () => {
        addRequestToList('Test request');

        const list = document.getElementById('requestList');
        expect(list.children.length).toBe(1);
        expect(list.firstChild.textContent).toBe('Test request: ');
        expect(list.firstChild.className).toBe('request-detail');
    });

    // Add more tests as needed
});