// Fetch posts from the API and render them
export async function fetchPosts(category = "all") {
  try {
    const url = category === "all" ? "/posts" : `/posts?category=${encodeURIComponent(category)}`;
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }

    const posts = await response.json();

    const postsContainer = document.getElementById("posts");
    postsContainer.innerHTML = ""; // Clear any existing content

    if (posts.length === 0) {
      postsContainer.innerHTML = "<p>No posts found.</p>";
      return;
    }
    // handling logout/login button visibility
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
      const postCard = createPostCard(post);
      postsContainer.appendChild(postCard);
    });
    if (posts[1] === 0) {
      document.querySelectorAll(".comment-input, .send-comment-btn").forEach((elem) => {
          elem.style.display = "none";
        });
    }
  } catch (error) {
    console.error("Error fetching posts:", error);
    const postsContainer = document.getElementById("posts");
    postsContainer.innerHTML = `<p>Error loading posts: ${error.message}</p>`;
  }
}

// Function to create a post card element
function createPostCard(post) {
  const postCard = document.createElement("div");
  postCard.className = "post-card";

  postCard.innerHTML = `
    <div class="title">${escapeHTML(post.Title)}</div>
    <div class="post-username">by @${escapeHTML(post.UserName)}</div>
    <div class="post-content">${escapeHTML(post.Content)}</div>
    <div class="details-toggle" onclick="toggleDetails(this)">
      <span class="details-text">Details</span>
    </div>
    <div class="meta hidden">${escapeHTML(post.Category)}, ${timeAgo(post.CreatedAt)}</div>
    <div class="post-actions">
      <button class="post-btn like" style="background:none;" id="${post.ID}">❤️</button>
      <div class="post-likes like">${escapeHTML(post.Likes.toString())}</div>
      <button class="post-btn dislike" style="background:none;" id="${post.ID}">👎</button>
      <div class="post-dislikes">${escapeHTML(post.Dislikes.toString())}</div>
    </div>
    <button class="comment-btn" onclick="toggleComments(${post.ID}, this)">Show Comments</button>
    <div class="comment-section hidden" id="comment-section-${post.ID}">
      <textarea class="comment-input" id="comment-input-${post.ID}" placeholder="Your comment"></textarea>
      <button class="send-comment-btn" onclick="postComment(${post.ID}, 1)">Comment</button>
      <div id="comments-list-${post.ID}" class="comments-list"></div>
    </div>
  `;
  
  likeEvent(postCard);
  return postCard;
}

window.fetchPosts = fetchPosts;
