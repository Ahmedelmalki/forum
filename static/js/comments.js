/*************Start Comment sections functions*****************/
// Function to post a comment
async function postComment(postId, userId) {
  const commentInput = document.getElementById(`comment-input-${postId}`);
  const commentContent = commentInput.value;

  try {
    const response = await fetch("/comments", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        post_id: postId,
        user_id: userId,
        content: commentContent,
      }),
    });

    if (response.ok) {
      loadComments(postId);
      commentInput.value = "";
    }
  } catch (error) {
    console.error("Error of posting comment:", error);
  }
}

// Function to load comments
async function loadComments(postId) {
  try {
    const response = await fetch(`/comments?post_id=${postId}`);
    const comments = await response.json();

    const commentsList = document.getElementById(`comments-list-${postId}`);
    commentsList.innerHTML = "";

    comments.reverse().forEach((comment) => {
      const commentElement = document.createElement("div");
      function timeAgo(date) {
        const seconds = Math.floor((new Date() - new Date(date)) / 1000);
        const intervals = [
          { label: "year", seconds: 31536000 },
          { label: "month", seconds: 2592000 },
          { label: "day", seconds: 86400 },
          { label: "hour", seconds: 3600 },
          { label: "minute", seconds: 60 },
          { label: "second", seconds: 1 },
        ];

        for (const interval of intervals) {
          const count = Math.floor(seconds / interval.seconds);
          if (count > 0) {
            return `${count} ${interval.label}${count !== 1 ? "s" : ""} ago`;
          }
        }
        return "just now";
      }

      commentElement.innerHTML = `
        <div class="comment">
          <small>Posted by <b>@${comment.username}</b>, ${timeAgo(
          comment.created_at
          )}</small>
          <p>${escapeHTML(comment.content)}</p>
          <div class="comment-actions">
            <button class="comment-btn like" style="background:none;" id="${comment.ID}">❤️</button>
            <span class="comment-likes like">${comment.Likes.toString()}</span>
            <button class="comment-btn dislike", style="background:none;"  id = ${comment.ID}>👎</button>
            <span class="comment-dislikes dislike">${comment.DislikeCount.}</span>
          </div>
        </div>
      `;
      likeCommentEvent(commentElement);
      commentsList.appendChild(commentElement);
    });
  } catch (error) {
    console.error("RError of loading comments:", error);
  }
}
/*************End Comment sections functions*****************/

// Load posts when the page loads
//window.onload = fetchPosts;

/*********************Likes on comments******************/
async function UpdateLikeOnComment(comment) {
  try {
    const response = await fetch("/comment-reactions");
    console.log("Fetching done");

    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }
    const likes = await response.json();
    console.log("LIkes fetched successfully");
    comment.querySelector(
      ".comment-actions .comment-likes"
    ).textContent = `${likes.LikeCount} likes`;
    comment.querySelector(
      ".comment-actions .comment-dislikes"
    ).textContent = `${likes.DislikeCount} dislikes`;
  } catch (err) {
    console.error("Error fetching likes:", err);
  }
}
function likeCommentEvent(comment) {
  likeButton = comment.querySelectorAll(".comment-actions .comment-btn");
  console.log(likeButton);

  if (window.cookie == "") {
    likeButton.disabled = true;
    likeButton.style.backgroundcolor = "#a9a9a9";
    likeButton.style.cursor = "not-allowed";
  } else {
    likeButton.forEach((element) => {
      element.addEventListener("click", async () => {
        try {
          // console.log( likeButton.classList.contains("like"))
          const response = await fetch("/comment-reactions", {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({
              UserId: 0,
              CommentId: parseInt(
                comment.querySelector(".comment-actions .comment-btn").id
              ),
              LikeCount: 0,
              Type: element.classList.contains("dislike") ? "dislike" : "like",
            }),
          });
          if (!response.ok) {
            const err = document.querySelector(".error-mssg");
            if (!err) {
              const erroemssg = document.createElement("p");
              erroemssg.className = "error-mssg";
              erroemssg.innerHTML = "user not found";
              document.appendChild(erroemssg);
            }
          }
          await UpdateLikeOnComment(comment);
        } catch (err) {
          console.log(err);
        }
      });
    });
  }
}