document.addEventListener("DOMContentLoaded", function () {
  const commentsContainer = document.getElementById("comments-container");
  const loadMoreBtn = document.getElementById("load-more-btn");
  const comments = commentsContainer.getElementsByClassName("comment-item");
  const commentsPerPage = 10;
  let currentlyShown = commentsPerPage;

  // Initially hide comments beyond the first 10
  for (let i = commentsPerPage; i < comments.length; i++) {
    comments[i].style.display = "none";
  }

  // Hide the "Load more" button if there are 10 or fewer comments
  if (comments.length <= commentsPerPage) {
    loadMoreBtn.style.display = "none";
  }

  loadMoreBtn.addEventListener("click", function () {
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
  });
});
