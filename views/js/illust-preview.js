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

  const image = url.replace(/c\/\d+x\d+.*?\//, "").replace(/square1200/, "master1200");
  const img = document.createElement("img");
  img.src = image;
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

function Overlay() {
  document.querySelectorAll('.artwork-small img').forEach(illust => {
    const url = illust.getAttribute("src");
    const button = document.createElement('div');
	button.style.cssText = `
	  position: absolute;
	  top: 0;
	  left: 0;
	  width: 100%;
	  height: 100%;
		    `;
	illust.parentElement.parentElement.appendChild(button);
	button.onclick = (e) => {
	  OpenPreviewer(url);
	};
      })
  console.log("ran")
}

Overlay();
