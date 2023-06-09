// polyfill
if (
  typeof globalThis.chrome !== "undefined" &&
  (typeof globalThis.browser === "undefined" ||
    Object.getPrototypeOf(globalThis.browser) !== Object.prototype)
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

export function getServer() {
  return browser.storage.local.get("server").then((x) => x.server);
}

export async function getClientID() {
  const q = await browser.storage.local.get("client_id");
  if (typeof q.client_id === "undefined") {
    const client_id = crypto.randomUUID();
    await browser.storage.local.set({ "client_id": client_id });
    return client_id;
  }
  return q.client_id;
}

export function humanTimeDiff(seconds) {
  const secInMin = 60;
  const minInHour = 60;
  const hourInDay = 24;
  let val = seconds;
  let label = "second";
  if (seconds >= secInMin * minInHour * hourInDay) {
    val = Math.floor(seconds / (secInMin * minInHour * hourInDay));
    label = "day";
  } else if (seconds >= secInMin * minInHour) {
    val = Math.floor(seconds / (secInMin * minInHour));
    label = "hour";
  } else if (seconds >= secInMin) {
    val = Math.floor(seconds / secInMin);
    label = "minute";
  }
  if (val == 1) {
    return `${val} ${label}`;
  } else {
    return `${val} ${label}s`;
  }
}
