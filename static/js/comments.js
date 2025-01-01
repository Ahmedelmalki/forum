async function postComment(postId, userId) {
  const commentInput = document.getElementById(`comment-input-${postId}`);
  const commentContent = commentInput.value;
  console.log("🚨🚨🚨🚨🚨🚨🚨🚨",commentInput.value)

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
    console.error("Error posting comment:", error);
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
      const commentElement = createCommentElement(comment);
      commentsList.appendChild(commentElement);
    });
  } catch (error) {
    console.error("Error loading comments:", error);
  }
}

// Function to create a comment element
function createCommentElement(comment) {
  console.log("============>>",comment)
  const commentElement = document.createElement("div");
  
  commentElement.innerHTML = `
    <div class="comment">
      <small>Posted by <b>@${comment.username}</b>, ${timeAgo(comment.created_at)}</small>
      <p>${escapeHTML(comment.content)}</p>
      <div class="comment-actions">
      <button class="like-btn">👍🏽</button>
      <button class="like-btn">👎🏽</button>
      </div>
    </div>
  `;
  likeEventOnComment(commentElement)
  return commentElement;
}
