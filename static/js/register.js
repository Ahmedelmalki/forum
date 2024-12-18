// helper func
function validateCredentials(email, password) {
  const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
  /* At least one letter,
  At least one digit, At least one special character
  Minimum of 8 characters */
  const passwordRegex =
    /^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/;

  const isEmailValid = emailRegex.test(email);
  const isPasswordStrong = passwordRegex.test(password);

  if (!isEmailValid) {
    alert("Invalid email format.");
    return false;
  }

  if (!isPasswordStrong) {
    alert(
      "Password must be at least 8 characters long, include at least one letter, one number, and one special character."
    );
    return false;
  }

  alert("Credentials are valid.");
  return true;
}

const form = document.querySelector(".container");
const botton = document.querySelector(".submit");
botton.addEventListener("click", async () => {
  const usernameInput = document.getElementById("username");
  const emailInput = document.getElementById("email");
  const passwordInput = document.getElementById("password");

  if (!validateCredentials(emailInput.value, passwordInput.value)) {
    console.log("here");
    return; // Stop if validation fails
  }
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
