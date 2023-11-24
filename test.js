const {addRequestToList, fetchInitialData, clearRequests} = require('./static/assets/script');

test('addRequestToList adds a new item to the list', () => {
    document.body.innerHTML = '<ul id="requestList"></ul>';

    addRequestToList('Test request');

    const list = document.getElementById('requestList');
    expect(list.children.length).toBe(1);
    expect(list.firstChild.textContent).toBe('Test request: ');
    expect(list.firstChild.className).toBe('request-detail');
});


test('fetchInitialData data sends request and populates list with response', async () => {
    global.fetch = jest.fn(() =>
        Promise.resolve({
            ok: true,
            json: () => Promise.resolve(["request1", "request2"])
        })
    );

    document.body.innerHTML = '<ul id="requestList"></ul>';

    await fetchInitialData();

    // allow microtasks to finish
    setTimeout(() => {
        const list = document.getElementById('requestList');
        expect(list.children.length).toBe(2);
        expect(list.firstChild.textContent).toBe('request2:');
        expect(list.firstChild.className).toBe('request-detail');
        expect(list.children[1].textContent).toBe('request1:');
        expect(list.children[1].className).toBe('request-detail');
    }, 0);

    expect(global.fetch).toHaveBeenCalled();
    expect(global.fetch).toHaveBeenCalledWith('/api/requests/all');

    jest.restoreAllMocks();
});


test('clearRequests clears list and sends request', async () => {
    global.fetch = jest.fn(() =>
        Promise.resolve({ok: true, json: () => Promise.resolve()})
    );

    document.body.innerHTML = '<ul id="requestList"><li>request 1</li></ul>';

    await clearRequests();

    const list = document.getElementById('requestList');
    expect(list.children.length).toBe(0);
    expect(list.innerHTML).toBe('');

    expect(global.fetch).toHaveBeenCalled();
    expect(global.fetch).toHaveBeenCalledWith('/api/requests/delete', {method: 'POST'});

    jest.restoreAllMocks();
});