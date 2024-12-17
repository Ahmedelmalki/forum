// Fetch posts from the API and render them
async function fetchPosts() {
  try {
    const response = await fetch("/posts");
    console.log("Fetching done");

    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }

    const posts = await response.json();
    console.log("Posts fetched successfully");

    const postsContainer = document.getElementById("posts");
    postsContainer.innerHTML = ""; // Clear any existing content

    if (posts.length === 0) {
      postsContainer.innerHTML = "<p>No posts found.</p>";
      return;
    }

    posts.forEach((post) => {
      console.log("Rendering post");

      const postCard = document.createElement("div");
      postCard.className = "post-card";

      postCard.innerHTML = `
      <div class="post-title">${escapeHTML(post.Title)}</div>
      <div class="meta">
        Category: ${escapeHTML(post.Category)} | 
        Posted on: ${new Date(post.CreatedAt).toLocaleString()}
      </div>
      <div class="post-content">${escapeHTML(post.Content)}</div>
      <div class="post-actions">
        <button class="post-btn">Like</button>
        <button class="comment-btn" onclick="loadComments(${post.ID})">
          View Comments
        </button>
        <div class="post-likes">${post.Likes || 0} likes</div>
      </div>
      <div class="comment-section">
        <textarea id="comment-input-${post.ID}" placeholder="Your comment"></textarea>
        <button class="send-comment-btn" onclick="postComment(${post.ID}, 1)">Comment</button>
        <div id="comments-list-${post.ID}" class="comments-list"></div>
      </div>
    `;
      postsContainer.appendChild(postCard);
    });
  } catch (error) {
    console.error("Error fetching posts:", error);
    const postsContainer = document.getElementById("posts");
    postsContainer.innerHTML = `<p>Error loading posts: ${error.message}</p>`;
  }
}

// Utility function to escape HTML to prevent XSS
function escapeHTML(str) {
  if (typeof str !== "string") return "";
  return str.replace(
    /[&<>'"]/g,
    (tag) =>
      ({
        "&": "&amp;",
        "<": "&lt;",
        ">": "&gt;",
        "'": "&#39;",
        '"': "&quot;",
      }[tag] || tag)
  );
}

/*************Start Comment sections functions*****************/
// Function to post a comment
async function postComment(postId, userId) {
  const commentInput = document.getElementById(`comment-input-${postId}`);
  const commentContent = commentInput.value;
  
  try {
    const response = await fetch('/comments', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        post_id: postId,
        user_id: userId,
        content: commentContent
      })
    });

    if (response.ok) {
      loadComments(postId);
      commentInput.value = '';
    }
  } catch (error) {
    console.error('Error of posting comment:', error);
  }
}

// Function to load comments
async function loadComments(postId) {
  try {
    const response = await fetch(`/comments?post_id=${postId}`);
    const comments = await response.json();
    
    const commentsList = document.getElementById(`comments-list-${postId}`);
    commentsList.innerHTML = ''; 

    comments.reverse().forEach(comment => {
      const commentElement = document.createElement('div');
      commentElement.innerHTML = `
      <div class="comment">
        <small>Posted by ${comment.username}, at: ${new Date(comment.created_at).toLocaleString()}</small>
        <p>${escapeHTML(comment.content)}</p>
        <button class="like-btn" onclick="likeComment(${comment.id})">Like</button>
        <button class="delete-btn" onclick="deleteComment(${comment.id})">Unlike</button>
      </div>
      `;
      commentsList.appendChild(commentElement);
    });
  } catch (error) {
    console.error('RError of loading comments:', error);
  }
}
/*************End Comment sections functions*****************/

// Load posts when the page loads
window.onload = fetchPosts;
