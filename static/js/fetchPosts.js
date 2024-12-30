// Fetch posts from the API and render them
export async function fetchPosts(category = "all") {
  try {
    const url =
      category === "all"
        ? "/posts"
        : `/posts?category=${encodeURIComponent(category)}`;
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
         <div class="post-username">by @${escapeHTML(post.UserName)}</div>
         
         <div class="post-content">${escapeHTML(post.Content)}</div>
        <div class="details-toggle" onclick="toggleDetails(this)">
           <span class="details-text">Details</span>
        </div>
         <div class="meta hidden">
          ${escapeHTML(post.Category)}, ${timeAgo(post.CreatedAt).toLocaleString()}
        </div>
         <div class="post-actions">
         <button class="post-btn like" id="like-${post.ID}" onclick="handleLike(${post.ID})">❤️</button>
         <div class="post-likes">${escapeHTML(post.Likes.toString())}</div>

         <button class="post-btn dislike" id="dislike-${post.ID}" onclick="handleDislike(${post.ID})">👎</button>
         <div class="post-dislikes">${escapeHTML(post.Dislikes.toString())}</div>
         </div>

         <button class="comment-btn" onclick="toggleComments(${post.ID}, this)">Show Comments</button>
        <div class="comment-section hidden" id="comment-section-${post.ID}">
          <textarea class="comment-input" id="comment-input-${post.ID}" placeholder="Your comment"></textarea>
          <button class="send-comment-btn" onclick="postComment(${post.ID}, 1)">Comment</button>
          <div id="comments-list-${post.ID}" class="comments-list"></div>

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

// entry point
// document.addEventListener("DOMContentLoaded", () => {
//   //document.getElementById("apply-filter").click();
//   fetchPosts('all');
//   // filtring logic
//   document.getElementById("apply-filter").addEventListener("click", () => {
//     const category = document.getElementById("category-filter").value;
//     fetchPosts(category);
//   });
// });

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
    loadComments(postId); // Fetch and display comments
  } else {
    console.log("Hiding comments for post:", postId);
    commentSection.classList.add("hidden");
    button.textContent = "Show Comments";
  }
}

export function toggleDetails(toggleElement) {
  const meta = toggleElement.nextElementSibling; // Select the `.meta` div
  meta.classList.toggle("hidden"); // Toggle the `hidden` class

  const detailsText = toggleElement.querySelector(".details-text");
  detailsText.textContent = meta.classList.contains("hidden")
    ? "Details"
    : "Hide Details";
}
// Add event listener for filter toggle
document.addEventListener("DOMContentLoaded", () => {
  const filterToggle = document.getElementById("filter-toggle");
  const filterForm = document.getElementById("filter-form");

  filterToggle.addEventListener("click", () => {
    filterForm.classList.toggle("hidden");
  });

  // Filtering logic
  document.getElementById("apply-filter").addEventListener("click", (e) => {
    e.preventDefault(); // Prevent default form submission
    const category = document.getElementById("category-filter").value;
    fetchPosts(category); // Call fetchPosts with the selected category
  });
});


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

//window.onload = fetchPosts;
window.escapeHTML = escapeHTML;
window.fetchPosts = fetchPosts;
window.toggleComments = toggleComments;
window.toggleDetails = toggleDetails;