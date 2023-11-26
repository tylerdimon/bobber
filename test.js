const {addRequestToList, clearRequests} = require('./static/assets/script');

test('addRequestToList adds a new item to the list', () => {
    document.body.innerHTML = '<ul id="requestList"><li>list item</li></ul>';

    addRequestToList('<li>another list item</li>');

    const list = document.getElementById('requestList');
    expect(list.children.length).toBe(2);
    expect(list.firstChild.textContent).toBe('another list item');
    expect(list.children[1].textContent).toBe('list item');
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