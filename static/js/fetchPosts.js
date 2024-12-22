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
  <div class="post-actions">
    <button class="post-btn">Like</button>
    <button class="post-btn-dislike">Dislike</button>
    <button class="comment-btn" onclick="toggleComments(${post.ID}, this)">
      Show Comments
    </button>
    <div class="post-likes">${post.Likes || 0} likes</div>
  </div>
  <div class="details-toggle" onclick="toggleDetails(this)">
    <span class="details-text">Details</span>
  </div>
  <div class="meta hidden">
    ${escapeHTML(post.Category)}, ${timeAgo(post.CreatedAt).toLocaleString()}
  </div>

  <div class="comment-section hidden" id="comment-section-${post.ID}">
    <textarea class="comment-input" id="comment-input-${post.ID}" placeholder="Your comment"></textarea>
    <button class="send-comment-btn" onclick="postComment(${post.ID}, 1)">Comment</button>
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
//sound handiling
// Select audio and button elements
const audio = document.getElementById("background-audio");
const muteButton = document.getElementById("mute-button");
const muteIcon = document.getElementById("mute-icon");

// Define SVG paths for mute and unmute states
const muteSVG = `
    <path d="M12 3.5L6 8H2v8h4l6 4.5V3.5zM16.5 12l2.8-2.8c.4-.4.4-1 0-1.4s-1-.4-1.4 0l-2.8 2.8-2.8-2.8c-.4-.4-1-.4-1.4 0s-.4 1 0 1.4l2.8 2.8-2.8 2.8c-.4.4-.4 1 0 1.4s1 .4 1.4 0l2.8-2.8 2.8 2.8c.4.4 1 .4 1.4 0s.4-1 0-1.4L16.5 12z" />
`;
const unmuteSVG = `
    <path d="M12 3.5L6 8H2v8h4l6 4.5V3.5zM18 6c-.6 0-1 .4-1 1v10c0 .6.4 1 1 1s1-.4 1-1V7c0-.6-.4-1-1-1zM15 9c-.6 0-1 .4-1 1v4c0 .6.4 1 1 1s1-.4 1-1v-4c0-.6-.4-1-1-1z" />
`;


// Toggle mute state and update SVG
muteButton.addEventListener("click", () => {
    if (audio.muted) {
        // Unmute
        audio.muted = false;
        muteIcon.innerHTML = muteSVG; // Update to mute icon
        muteButton.setAttribute("aria-label", "Mute");
    } else {
        // Mute
        audio.muted = true;
        muteIcon.innerHTML = unmuteSVG; // Update to unmute icon
        muteButton.setAttribute("aria-label", "Unmute");
    }
});

// Optional: Debugging information
console.log("Audio element:", audio);
console.log("Mute button:", muteButton);


// Pause audio when the tab is not visible
document.addEventListener("visibilitychange", () => {
    if (document.visibilityState === "hidden") {
        audio.pause();
    } else if (document.visibilityState === "visible") {
        audio.play();
    }
});

// Define start and end times in seconds
const startTime = 120;
const endTime = 146;

// Set the audio to start playing from the specified start time
audio.currentTime = startTime;

// Loop the specific section of the audio
audio.addEventListener("timeupdate", () => {
    if (audio.currentTime >= endTime) {
        audio.currentTime = startTime;
        audio.play();
    }
});

function toggleComments(postId, button) {
  const commentSection = document.getElementById(`comment-section-${postId}`);
  console.log("Button clicked:", button.textContent);
  console.log("Comment section hidden:", commentSection.classList.contains('hidden'));

  if (commentSection.classList.contains('hidden')) {
    console.log("Showing comments for post:", postId);
    commentSection.classList.remove('hidden');
    button.textContent = 'Hide Comments';
    loadComments(postId); // Fetch and display comments
  } else {
    console.log("Hiding comments for post:", postId);
    commentSection.classList.add('hidden');
    button.textContent = 'Show Comments';
  }
}

// function likePost(postId) {
//   console.log(`Post ${postId} liked!`);
//   // Implement API call to like the post and update UI accordingly
// }

// function unlikePost(postId) {
//   console.log(`Post ${postId} unliked!`);
//   // Implement API call to unlike the post and update UI accordingly
// }


function toggleDetails(toggleElement) {
  const meta = toggleElement.nextElementSibling; // Select the `.meta` div
  meta.classList.toggle('hidden'); // Toggle the `hidden` class

  const detailsText = toggleElement.querySelector('.details-text');
  detailsText.textContent = meta.classList.contains('hidden') ? 'Details' : 'Hide Details';
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
