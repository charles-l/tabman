var tabst = document.getElementById('tabs');
document.getElementById('settings').addEventListener('click', () => browser.runtime.openOptionsPage());
async function main() {
    var serverName = await getServer();
    if(serverName === undefined) {
        tabst.innerHTML = "<b>no server set -- please configure the TabMan server in the settings</b>"
        return;
    }
    console.log('server', serverName);
    var response = await fetch(`http://${serverName}/tabs/`, {
        headers: {
            "Content-Type": "application/json",
        },
    });

    var j = await response.json();
    tabst.innerHTML =
        j.map(instance =>
            "<b>" + instance.client_id + "</b>" +
            "<table>" + instance.tabs.map(r =>
                `<tr><td><img src="${r[0]}" width="15px"></img><a href="${r[2]}">${r[1]}</a></td></tr>`).join("") + "</table>"
        );
}
main()
