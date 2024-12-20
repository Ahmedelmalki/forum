function validateCredentials(email, password) {
  const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
  const passwordRegex =
    /^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/;

  const isEmailValid = emailRegex.test(email);
  const isPasswordStrong = passwordRegex.test(password);

  const errors = [];
  if (!isEmailValid) errors.push("Invalid email format");
  if (!isPasswordStrong)
    errors.push(
      "Password must be at least 8 characters long, include at least one letter, one number, and one special character"
    );

  return {
    isValid: isEmailValid && isPasswordStrong,
    errors: errors,
  };
}

const button = document.querySelector(".submit");
button.addEventListener("click", async () => {
  const validation = validateCredentials(emailInput.value, passwordInput.value);

  if (!validation.isValid) {
    // Display errors in the UI instead of using alerts
    const errorDiv = document.createElement("div");
    errorDiv.className = "error-messages";
    validation.errors.forEach((error) => {
      const p = document.createElement("p");
      p.textContent = error;
      errorDiv.appendChild(p);
    });
    form.appendChild(errorDiv);
    return;
  }

  // proceed with form submission...
});

const form = document.querySelector(".container");
button.addEventListener("click", async () => {
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
