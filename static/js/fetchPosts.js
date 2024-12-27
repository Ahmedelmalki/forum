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

      console.log(post.Likes, typeof post.Likes);
      postCard.innerHTML = `
        <div class="post-username">${escapeHTML(post.UserName)}</div>
        <div class="meta">
        Title: ${escapeHTML(post.Title)} |
          Category: ${escapeHTML(post.Category)} | 
          Posted on: ${new Date(post.CreatedAt).toLocaleString()}
        </div>
        <div class="post-content">${escapeHTML(post.Content)}</div>
        <div class="post-actions">
          <button class="post-btn like", id = ${post.ID}>Like</button>
          <div class="post-likes like">${escapeHTML(
            post.Likes.toString()
          )} likes</div>
                  <button class="post-btn dislike", style = "background:crimson",  id = ${
                    post.ID
                  }>Dislike</button>
          <div class="post-likes " >${escapeHTML(
            post.Likes.toString()
          )} dislikes</div>

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
      likeEvent(postCard);

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

/********************************** likes ******************************* */

async function UpdateLike(post) {
  try {
    const response = await fetch("/like");
    console.log("Fetching done #### here");

    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }
    const likes = await response.json();
    console.log("LIkes fetched successfully");
    post.querySelector(
      ".post-actions .post-likes"
    ).textContent = `${likes.LikeCOunt} likes`;
  } catch (err) {
    console.error("Error fetching likes:", err);
  }
}
function likeEvent(post) {
  likeButton = post.querySelector(".post-actions .post-btn");
  if (document.cookie == "") {
    likeButton.disabled = true;
    likeButton.style.background = "#a9a9a9";
    likeButton.style.cursor = "not-allowed";
  } else {
    likeButton.addEventListener("click", async () => {
      try {
        console.log(likeButton.classList.contains("like"));
        const response = await fetch("/like", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            UserId: 55,
            PostId: parseInt(post.querySelector(".post-actions .post-btn").id),
            LikeCOunt: 4,
            Type: likeButton.classList.contains("dislike") ? "dislike" : "like",
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
        await UpdateLike(post);
      } catch (err) {
        console.log(err);
      }
    });
  }
}
