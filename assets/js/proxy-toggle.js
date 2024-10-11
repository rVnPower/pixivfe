document.addEventListener("DOMContentLoaded", function () {
  const showAllProxiesToggle = document.getElementById("show-all-proxies");
  const imageProxySelect = document.getElementById("image-proxy");
  const proxyCheckEnabled = imageProxySelect.getAttribute("data-proxy-check-enabled") === "true";

  function updateProxyVisibility() {
    const showAll = showAllProxiesToggle.checked;
    const options = imageProxySelect.options;
    let firstVisibleOption = null;

    for (let i = 0; i < options.length; i++) {
      const option = options[i];
      const proxyType = option.getAttribute("data-proxy-type");

      if (proxyType === "working" || (proxyType === "all" && showAll) || !proxyCheckEnabled) {
        option.style.display = "";
        if (!firstVisibleOption) firstVisibleOption = option;
      } else {
        option.style.display = "none";
      }
    }

    // If the currently selected option is now hidden, select the first visible option
    if (imageProxySelect.selectedOptions[0].style.display === "none") {
      firstVisibleOption.selected = true;
    }
  }

  // Set initial state
  if (!proxyCheckEnabled) {
    showAllProxiesToggle.checked = true;
    showAllProxiesToggle.disabled = true;
  }

  updateProxyVisibility();

  showAllProxiesToggle.addEventListener("change", updateProxyVisibility);
});
