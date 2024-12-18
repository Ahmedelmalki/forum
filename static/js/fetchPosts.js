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

    if (posts[1] != 0) {
      document.querySelectorAll(".loged").forEach((elem) => {
        elem.style.display = "none";
      });
    }
    if (posts[1] === 0) {
      document.querySelectorAll(".unloged").forEach((elem) => {
        elem.style.display = "none";
      });
    }

    posts[0].forEach((post) => {
      console.log("Rendering post");

      const postCard = document.createElement("div");
      postCard.className = "post-card";

      postCard.innerHTML = `
      <div class="post-username">${escapeHTML(post.UserName)}</div>
      <div class="meta">
        Title: ${escapeHTML(post.Title)} |
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
        <textarea class="comment-input" id="comment-input-${
          post.ID
        }" placeholder="Your comment"></textarea>
        <button class="send-comment-btn" onclick="postComment(${
          post.ID
        }, 1)">Comment</button>
        <div id="comments-list-${post.ID}" class="comments-list"></div>
      </div>
    `;
      postsContainer.appendChild(postCard);
    });
    if (posts[1] === 0) {
      document
        .querySelectorAll(".comment-input, .send-comment-btn")
        .forEach((elem) => {
          elem.style.display = "none";
        });
    }
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

window.onload = fetchPosts;
