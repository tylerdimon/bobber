const {addRequestToList, deleteRequests} = require('./static/assets/script');

test('addRequestToList adds a new item to the list', () => {
    document.body.innerHTML = '<ul id="request-list"><li>list item</li></ul>';

    addRequestToList('<li>another list item</li>');

    const list = document.getElementById('request-list');
    expect(list.children.length).toBe(2);
    expect(list.firstChild.textContent).toBe('another list item');
    expect(list.children[1].textContent).toBe('list item');
});


test('deleteRequests clears list and sends request to API', async () => {
    global.fetch = jest.fn(() =>
        Promise.resolve({ok: true, json: () => Promise.resolve()})
    );

    document.body.innerHTML = '<ul id="request-list"><li>request 1</li></ul>';

    await deleteRequests();

    const list = document.getElementById('request-list');
    expect(list.children.length).toBe(0);
    expect(list.innerHTML).toBe('');

    expect(global.fetch).toHaveBeenCalled();
    expect(global.fetch).toHaveBeenCalledWith('/api/requests/delete', {method: 'POST'});

    jest.restoreAllMocks();
});