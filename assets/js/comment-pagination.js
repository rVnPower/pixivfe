let commentsContainer;
let loadMoreBtn;
let comments;
const commentsPerPage = 10;
let currentlyShown;
let isInitialized = false;

function initializeCommentPagination() {
  // Reset state variables
  commentsContainer = document.getElementById("comments-container");
  loadMoreBtn = document.getElementById("load-more-btn");
  comments = commentsContainer ? commentsContainer.getElementsByClassName("comment-item") : [];
  currentlyShown = 0;
  isInitialized = false;

  if (!commentsContainer) {
    // No comments container found, pagination not needed
    if (loadMoreBtn) loadMoreBtn.style.display = "none";
    return;
  }

  if (comments.length === 0) {
    // No comments, hide the "Load more" button
    if (loadMoreBtn) loadMoreBtn.style.display = "none";
    return;
  }

  if (comments.length <= commentsPerPage) {
    // 10 or fewer comments, no pagination needed
    Array.from(comments).forEach(comment => comment.style.display = "block");
    if (loadMoreBtn) loadMoreBtn.style.display = "none";
    return;
  }

  // More than 10 comments, initialize pagination
  currentlyShown = commentsPerPage;

  Array.from(comments).forEach((comment, index) => {
    comment.style.display = index < commentsPerPage ? "block" : "none";
  });

  if (loadMoreBtn) {
    loadMoreBtn.style.display = "block";

    // Add click event listener for load-more button
    loadMoreBtn.addEventListener("click", function() {
      for (
        let i = currentlyShown;
        i < currentlyShown + commentsPerPage && i < comments.length;
        i++
      ) {
        comments[i].style.display = "block";
      }

      currentlyShown += commentsPerPage;

      if (currentlyShown >= comments.length) {
        loadMoreBtn.style.display = "none";
      }
    });
  }

  isInitialized = true;
}

// Initialize on DOMContentLoaded
document.addEventListener("DOMContentLoaded", initializeCommentPagination);

// Reinitialize on htmx:load event
document.addEventListener("htmx:load", initializeCommentPagination);
