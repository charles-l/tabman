function saveOptions() {
  const serverInput = document.querySelector("#server");
  browser.storage.local.set({
    server: serverInput.value,
  });
}

document.querySelector("#save").addEventListener("click", () => saveOptions());
document.addEventListener("DOMContentLoaded", () => {
  getServer().then((s) => {
    if (s !== undefined) document.querySelector("#server").value = s;
  });
});
