const tabst = document.getElementById("tabs");
document.getElementById("settings").addEventListener(
  "click",
  () => browser.runtime.openOptionsPage(),
);
async function main() {
  const serverName = await getServer();
  if (serverName === undefined) {
    tabst.innerHTML =
      "<b>no server set -- please configure the TabMan server in the settings</b>";
    return;
  }
  console.log("server", serverName);
  const response = await fetch(`http://${serverName}/tabs/`, {
    headers: {
      "Content-Type": "application/json",
    },
  });

  const j = await response.json();
  j.forEach(instance => {
    const table = document.createElement('table');
    table.innerHTML = instance.tabs.map((r) =>
      `<tr><td><img src="${r[0]}" width="15px"></img><a href="${r[2]}">${
        r[1]
      }</a></td></tr>`
    ).join("")
    const title = document.createElement('b');
    title.innerText = instance.client_id;
    tabst.appendChild(title);
    const a = document.createElement('a');
    a.innerText = 'delete';
    a.href = '#';
    a.addEventListener('click', () => {
      fetch(`http://${serverName}/tabs/${instance.client_id}`, {
        method: 'DELETE'
      }).then(() => window.location.reload());
    });

    tabst.appendChild(a);
    tabst.appendChild(table);
  });
}
main();
