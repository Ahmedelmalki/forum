// Fetch posts from the API and render them
export async function fetchPosts(type) {
  try {
    let posts = null;
    if (type === "posts" || type === "LikedPosts" || type === "CtreatedBy") {
      if (type == "LikedPosts" || type === "CtreatedBy") {
        const response = await fetch(`/${type}`, {
          method: "POST",
        });
        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }

        posts = await response.json();
      } else {
        const response = await fetch(`/${type}`);
        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }

        posts = await response.json();
      }
    } else {
      const currentUrl = window.location.href;
      const url = new URL(currentUrl);
      const params = url.searchParams;
      const category = params.getAll("categories");

      let link = `/${type}?`;
      category.forEach((element, index) => {
        link += `categories=` + element;
        if (index < category.length - 1) {
          link += "&";
        }
      });

      const response = await fetch(link, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (!response.ok) {
        throw new Error(`HTTP error! Status: ${response.status}`);
      }
      posts = await response.json();
    }

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
      const postCard = document.createElement("div");
      postCard.className = "post-card";

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
      postCard.innerHTML = `
         <div class="title">${escapeHTML(post.Title)}</div>
         <div class="post-username">by @${escapeHTML(post.Username)}</div>
         
         <div class="post-content">${escapeHTML(post.Content)}</div>
        <div class="details-toggle" onclick="toggleDetails(this)">
           <span class="details-text">Details</span>
        </div>
        <div class="meta hidden">
        ${post.Categories.join(", ")}, ${timeAgo(
        post.CreatedAt
      ).toLocaleString()}
        </div>
         <div class="post-actions">
          <button class="post-btn like" style="background:none;" id="${
            post.Id
          }">❤️</button>
          <div class="post-likes like">${escapeHTML(
            post.Likes.toString()
          )} </div>
          <button class="post-btn dislike", style="background:none;"  id = ${
            post.ID
          }>👎</button>
          <div class="post-dislikes" >${escapeHTML(
            post.Dislikes.toString()
          )} </div>
        </div>
         <button class="comment-btn" onclick="toggleComments(${
           post.Id
         }, this)">Show Comments</button>
        <div class="comment-section hidden" id="comment-section-${post.Id}">
          <textarea class="comment-input" id="comment-input-${
            post.Id
          }" placeholder="Your comment"></textarea>
          <button class="send-comment-btn" onclick="postComment(${
            post.Id
          }, 1)">Comment</button>
          <div id="comments-list-${post.Id}" class="comments-list"></div>

      `;
      likeEvent(postCard);
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

export function toggleComments(postId, button) {
  const commentSection = document.getElementById(`comment-section-${postId}`);
  console.log("Button clicked:", button.textContent);
  console.log(
    "Comment section hidden:",
    commentSection.classList.contains("hidden")
  );

  if (commentSection.classList.contains("hidden")) {
    console.log("Showing comments for post:", postId);
    commentSection.classList.remove("hidden");
    button.textContent = "Hide Comments";
    loadComments(postId);
  } else {
    console.log("Hiding comments for post:", postId);
    commentSection.classList.add("hidden");
    button.textContent = "Show Comments";
  }
}

export function toggleDetails(toggleElement) {
  const meta = toggleElement.nextElementSibling;
  meta.classList.toggle("hidden");

  const detailsText = toggleElement.querySelector(".details-text");
  detailsText.textContent = meta.classList.contains("hidden")
    ? "Details"
    : "Hide Details";
}

// Utility function to escape HTML to prevent XSS
export function escapeHTML(str) {
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

window.escapeHTML = escapeHTML;
window.fetchPosts = fetchPosts;
window.toggleComments = toggleComments;
window.toggleDetails = toggleDetails;
