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
          <button class="comment-btn like" style="background:none;">👍</button>
          <div class="comment-likes like">${(comment.Likes)} </div>
          <button class="comment-btn dislike" style="background:none;">👎</button>
          <div class="comment-dislikes">${(comment.Dislikes)}  </div>
              </div>
        </div>
      `;
      likeEvent(commentElement, comment.id , postId);
      commentsList.appendChild(commentElement);
    });
  } catch (error) {
    console.error("RError of loading comments:", error);
  }
}
