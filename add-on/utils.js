function getServer() {
    return browser.storage.sync.get().then(x => x.server)
}
