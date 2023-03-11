function postTabs() {
  browser.tabs.query({}).then(tabs => {
    const tabsJson = tabs.map(x => [x.title, x.url]);
    console.log(tabsJson);
    browser.runtime.sendNativeMessage('tabman', tabsJson);
  });
}

browser.tabs.onRemoved.addListener(postTabs);
browser.tabs.onCreated.addListener(postTabs);
