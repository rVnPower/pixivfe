const scrollableArea = document.getElementById('horizontal-scroll');

scrollableArea.addEventListener('wheel', function(event) {
  event.preventDefault(); // Prevent vertical scroll
  scrollableArea.scrollLeft += event.deltaY; // Scroll horizontally using deltaY
});
