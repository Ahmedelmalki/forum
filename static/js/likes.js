async function UpdateLike(post) {
  try {
    const response = await fetch("/like");
    console.log("Fetching done");

    if (!response.ok) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }
    const likes = await response.json();
    console.log("LIkes fetched successfully");
    post.querySelector(
      ".post-actions .post-likes"
    ).textContent = `${likes.LikeCOunt} likes`;
    post.querySelector(
      ".post-actions .post-dislikes"
    ).textContent = `${likes.DislikeCOunt} dislikes`;
  } catch (err) {
    console.error("Error fetching likes:", err);
  }
}
function likeEvent(post) {
  likeButton = post.querySelectorAll(".post-actions .post-btn");
  console.log(likeButton);

  if (window.cookie == "") {
    likeButton.disabled = true;
    likeButton.style.backgroundcolor = "#a9a9a9";
    likeButton.style.cursor = "not-allowed";
  } else {
    likeButton.forEach((element) => {
      element.addEventListener("click", async () => {
        try {
          // console.log( likeButton.classList.contains("like"))
          const response = await fetch("/like", {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({
              UserId: 0,
              PostId: parseInt(
                post.querySelector(".post-actions .post-btn").id
              ),
              LikeCOunt: 0,
              Type: element.classList.contains("dislike") ? "dislike" : "like",
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
    });
  }
}

// setInterval(()=>UpdateLike(post), 1000);
