// Fetch posts from the API and render them
async function fetchPosts() {
  try {
    const response = await fetch("/posts");
    console.log("feching done")

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const posts = await response.json();
    console.log("111111111111111")

    const postsContainer = document.getElementById("posts");
    postsContainer.innerHTML = ""; // Clear any existing content

    if (posts.length === 0) {
      postsContainer.innerHTML = "<p>No posts found.</p>";
      return;
    }

    posts.forEach((post) => {
      console.log("222222222222222222")
      const postDiv = document.createElement("div");
      postDiv.className = "post";

      postDiv.innerHTML = `
        <h2>${escapeHTML(post.Title)}</h2>
        <div class="meta">
          Category: ${escapeHTML(post.Category)} | 
          Posted on: ${new Date(post.CreatedAt).toLocaleString()}
        </div>
        <p>${escapeHTML(post.Content)}</p>
      `;
      postsContainer.appendChild(postDiv);
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
