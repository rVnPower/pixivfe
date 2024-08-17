function OpenPreviewer(url) {
  const viewer = document.createElement("div");
  viewer.style.cssText = `
    height: 100vh;
    width: 100vw;
    position: fixed;
    top: 0;
    left: 0;
    background: rgba(0,0,0,.8);
		display: flex;
		flex-direction: column;
		padding: 0 3rem;
		overflow: scroll;
  `;

  const imageLink = url.replace(/c\/\d+x\d+.*?\//, "").replace(/square1200/, "master1200");
  const img = document.createElement("img");
  img.src = imageLink;
  img.style.cssText = `
	margin: 3rem auto;
	max-width: 90%;
	max-height: 90%;
    `

  viewer.appendChild(img);

  document.body.appendChild(viewer);
  
  viewer.onclick = () => {
    document.body.removeChild(viewer);
  };
}

function AddOverlay() {
  // Check out `_layout.jet.html`
  const type = document.querySelector('#artworkPreview').innerHTML

  let className, html;

  if (type === "cover") {
    className = "overlay-cover";
    html = "";
  } else if (type === "button") {
    className = "overlay-button";
    html = "↗";
  } else {
    return;
  }

  document.querySelectorAll('.artwork-small .artwork-image img').forEach(illust => {
    const url = illust.getAttribute("src");
    const button = document.createElement('div');

    button.setAttribute("class", className);
    button.innerHTML = html;

    illust.parentElement.parentElement.appendChild(button);

    button.onclick = (e) => {
      OpenPreviewer(url);
    };
  })
}

addEventListener('htmx:afterSwap', function (event) {
  console.log("%o", event)
  AddOverlay();
});

// Initialize (it will only run one time)
AddOverlay();
