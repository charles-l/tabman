// polyfill
if(typeof globalThis.browser === "undefined" || Object.getPrototypeOf(globalThis.browser) !== Object.prototype) {
    // chrome
    globalThis.browser = globalThis.chrome;
    { // patch tabs.query
        const oldF = globalThis.browser.tabs.query;
        globalThis.browser.tabs.query = (query) => {
            return new Promise((cont) => {
                oldF(query, cont);
            });
        }
    }

    { // patch storage.local.get
        const oldLocal = globalThis.browser.storage.local
        globalThis.browser.storage.local = {
            get: (query) => {
                return new Promise((cont) => {
                    oldLocal.get(query, cont);
                });
            },
            set: (vals) => oldLocal.set(vals)
        }
    }
}

function getServer() {
    return browser.storage.local.get('server').then(x => x.server)
}

