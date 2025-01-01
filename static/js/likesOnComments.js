
// http.HandleFunc("/Commentlike", forum.HandleLikesOnComments(db))

/* <div class="post-actions">     post like shit
      <button class="post-btn like" style="background:none;" id="${post.ID}">❤️</button>
      <div class="post-likes like">${escapeHTML(post.Likes.toString())}</div>
      <button class="post-btn dislike" style="background:none;" id="${post.ID}">👎</button>
      <div class="post-dislikes">${escapeHTML(post.Dislikes.toString())}</div>
    </div>
*/

/*  <div class="comment">
      <small>Posted by <b>@${comment.username}</b>, ${timeAgo(comment.created_at)}</small>
      <p>${escapeHTML(comment.content)}</p>
      <div class="comment-actions">
      <button class="like-btn" id="">👍🏽</button>
      <div class="comment-likes">${escapeHTML(post.Likes.toString())}</div>
      <button class="like-btn" id="">👎🏽</button>
      <div class="comment-dislikes">${escapeHTML(post.Dislikes.toString())}</div>
      </div>
    </div>
*/

/*
type likesOnCmnts struct {
	User_Id      int    `json:"UserId"`
	Comment_Id   int    `json:"CommentId"`
	LikeCount    int    `json:"LikeCount"`
	DisLikeCount int    `json:"DislikeCount"`
	Type         string `json:"Type"`
}
*/
function likeEventOnComment(comment) {
  likeButton = comment.querySelectorAll(".comment-actions .like-btn");
  if (window.cookie == "") {
    likeButton.disabled = true;
    likeButton.style.backgroundcolor = "#a9a9a9";
    likeButton.style.cursor = "not-allowed";
  } else {
    likeButton.forEach((element) => {
      element.addEventListener("click", async () => {
        try {
          // console.log( likeButton.classList.contains("like"))
          const response = await fetch("/Commentlike", {
            
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({
              UserId: 0,
              CommentId: parseInt(comment.querySelector(".comment-actions .like-btn").id),
              LikeCOunt: 0,
              Type: element.classList.contains("dislike") ? "dislike" : "like",
            }),
          });
         // console.log(response);
          if (!response.ok) {
            const err = document.querySelector(".error-mssg");
            if (!err) {
              const erroemssg = document.createElement("p");
              erroemssg.className = "error-mssg";
              erroemssg.innerHTML = "user not found";
              document.appendChild(erroemssg);
            }
          }
          await UpdateLikeOnComment(comment);
        } catch (err) {
          console.log(err);
        }
      });
    });
  }
}

async function UpdateLikeOnComment(comment) {
  try {
    const response = await fetch("/Commentlike");
    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }
    const likesOnCmnts = await response.json();
    comment.querySelector(".comment-actions .comment-likes").textContent = `${likesOnCmnts.LikeCount} likes`;
    comment.querySelector(".comment-actions .comment-dislikes").textContent = `${likesOnCmnts.DislikeCount} dislikes`;
  } catch (err) {
    console.error("Error fetching likes:", err);
  }
}