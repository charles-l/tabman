// polyfill
if (
  typeof globalThis.browser === "undefined" ||
  Object.getPrototypeOf(globalThis.browser) !== Object.prototype
) {
  // chrome
  globalThis.browser = globalThis.chrome;
  { // patch tabs.query
    const oldF = globalThis.browser.tabs.query;
    globalThis.browser.tabs.query = (query) => {
      return new Promise((cont) => {
        oldF(query, cont);
      });
    };
  }

  { // patch storage.local.get
    const oldLocal = globalThis.browser.storage.local;
    globalThis.browser.storage.local = {
      get: (query) => new Promise((cont) => oldLocal.get(query, cont)),
      set: (vals) => new Promise((cont) => oldLocal.set(vals, cont)),
    };
  }
}

// deno-lint-ignore no-unused-vars
function getServer() {
  return browser.storage.local.get("server").then((x) => x.server);
}

// deno-lint-ignore no-unused-vars
async function getClientID() {
  const q = await browser.storage.local.get("client_id");
  if (typeof q.client_id === "undefined") {
    const client_id = crypto.randomUUID();
    await browser.storage.local.set({ "client_id": client_id });
    return client_id;
  }
  return q.client_id;
}
