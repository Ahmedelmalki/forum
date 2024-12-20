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
          <small>Posted by <b>@${comment.username}</b>, ${timeAgo(comment.created_at)}</small>
          <p>${escapeHTML(comment.content)}</p>
          <button class="like-btn" onclick="likeComment(${comment.id})">Like</button>
          <button class="delete-btn" onclick="deleteComment(${comment.id})">Unlike</button>
        </div>
      `;
      
      commentsList.appendChild(commentElement);
    });
  } catch (error) {
    console.error("RError of loading comments:", error);
  }
}
/*************End Comment sections functions*****************/

// Load posts when the page loads
window.onload = fetchPosts;
