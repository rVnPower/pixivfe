let commentsContainer;
let loadMoreBtn;
let comments;
const commentsPerPage = 10;
let currentlyShown;

function initializeCommentPagination() {
  commentsContainer = document.getElementById("comments-container");
  loadMoreBtn = document.getElementById("load-more-btn");

  if (!commentsContainer || !loadMoreBtn) return;

  comments = commentsContainer.getElementsByClassName("comment-item");
  currentlyShown = commentsPerPage;

  // Initially hide comments beyond the first 10
  for (let i = 0; i < comments.length; i++) {
    if (i < commentsPerPage) {
      comments[i].style.display = "block";
    } else {
      comments[i].style.display = "none";
    }
  }

  // Hide the "Load more" button if there are 10 or fewer comments
  if (comments.length <= commentsPerPage) {
    loadMoreBtn.style.display = "none";
  } else {
    loadMoreBtn.style.display = "block";
  }
}

document.addEventListener("DOMContentLoaded", initializeCommentPagination);
document.addEventListener("htmx:afterSwap", initializeCommentPagination);

document.addEventListener("click", function(event) {
  if (event.target && event.target.id === "load-more-btn") {
    // Show the next set of comments
    for (
      let i = currentlyShown;
      i < currentlyShown + commentsPerPage && i < comments.length;
      i++
    ) {
      comments[i].style.display = "block";
    }

    currentlyShown += commentsPerPage;

    // Hide the "Load more" button if all comments are shown
    if (currentlyShown >= comments.length) {
      loadMoreBtn.style.display = "none";
    }
  }
});
