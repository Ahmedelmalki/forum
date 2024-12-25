const form = document.querySelector(".container");
const botton = document.querySelector(".submit");
botton.addEventListener("click", async () => {
  const usernameInput = document.getElementById("username");
  const emailInput = document.getElementById("email");
  const passwordInput = document.getElementById("password");
  try {
    const response = await fetch("/register/submit", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        username: usernameInput.value,
        email: emailInput.value,
        password: passwordInput.value,
      }),
    });
    if (!response.ok) {
      const err = document.querySelector(".error-mssg");
      if (!err) {
        const errmssg = document.createElement("p");
        errmssg.className = "error-mssg";
        errmssg.innerHTML = "error registring please try again";
        form.appendChild(errmssg);
      }
    } else {
      window.location.href = "/";
    }
  } catch (err) {
    console.log(err);
  }
});
