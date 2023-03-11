function postTabs() {
  browser.tabs.query({}).then(tabs => {
    const tabsJson = tabs.map(x => [x.title, x.url]);
    console.log(tabsJson);
    fetch('http://localhost:8080/tabs/', {
      method: "POST",
      mode: "no-cors",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(tabsJson)
    }).then(x => console.log(x));
  });
}

browser.tabs.onUpdated.addListener(postTabs);
