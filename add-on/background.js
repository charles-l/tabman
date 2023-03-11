var favIconCache = {}

function blobToBase64(blob){
  return new Promise((cont) => {
    const reader = new FileReader();
    reader.readAsDataURL(blob);
    reader.onloadend = () => cont(reader.result);
  });
}

function postTabs() {
  browser.tabs.query({}).then(tabs => {
    // cache any icons that changed
    tabs.map(tab => {
      if(!(tab.favIconUrl in favIconCache)) {
        fetch(tab.favIconUrl).then(x => x.blob()).then(blobToBase64).then(x => favIconCache[tab.favIconUrl] = x).catch(err => favIconCache[tab.favIconUrl] = null);
      }
    })
    const tabsJson = {client_id: browser.runtime.id, tabs: tabs.map(x => [favIconCache[x.favIconUrl], x.title, x.url])};
    getServer().then(server => {
      if(server !== undefined) {
        fetch(`http://${server}/tabs/`, {
          method: "POST",
          headers: { "Content-Type": "application/json", },
          body: JSON.stringify(tabsJson)
        })
      } else {
        console.log("server is undefined, skipping post");
      }
    });
  });
}

browser.tabs.onUpdated.addListener(postTabs);
