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
          <div class="post-likes">${post.Likes || 0} likes</div>
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

// Load posts when the page loads
window.onload = fetchPosts;
