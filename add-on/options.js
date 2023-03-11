function saveOptions() {
    var serverInput = document.querySelector('#server');
    browser.storage.sync.set({
        server: serverInput.value
    })
}

document.querySelector('#server').addEventListener('change', () => saveOptions());
document.addEventListener('DOMContentLoaded', () => {
    getServer().then((s) => if(s !== undefined){document.querySelector('#server').value = s});
});
