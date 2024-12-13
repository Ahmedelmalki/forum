// Fetch posts from the API and render them
async function fetchPosts() {
  try {
    const response = await fetch("/api/posts"); // Fetch from the JSON endpoint
    const posts = await response.json();

    const postsContainer = document.getElementById("posts");
    postsContainer.innerHTML = ""; // Clear any existing content

    posts.forEach((post) => {
      const postDiv = document.createElement("div");
      postDiv.className = "post";

      postDiv.innerHTML = `
                        <h2>${post.title}</h2>
                        <div class="meta">Category: ${
                          post.category
                        } | Posted on: ${new Date(
        post.created_at
      ).toLocaleString()}</div>
                        <p>${post.content}</p>
                    `;
      postsContainer.appendChild(postDiv);
    });
  } catch (error) {
    console.error("Error fetching posts:", error);
  }
}

// Load posts when the page loads
window.onload = fetchPosts;
