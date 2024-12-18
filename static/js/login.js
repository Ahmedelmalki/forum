const form = document.querySelectorAll(".form-group")[0];
const botton = document.querySelector(".btn");
botton.addEventListener("click", async () => {
  const emailInput = document.getElementById("email");
  const passwordInput = document.getElementById("password");
  try {
    const response = await fetch("/login/submit", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        email: emailInput.value,
        password: passwordInput.value,
      }),
    });
    if (!response.ok) {
      const err = document.querySelector(".error-mssg");
      if (!err) {
        const erroemssg = document.createElement("p");
        erroemssg.className = "error-mssg";
        erroemssg.innerHTML = "user not found";
        form.appendChild(erroemssg);
      }
    } else {
      window.location.href = "/";
    }
  } catch (err) {
    console.log(err);
  }
});
