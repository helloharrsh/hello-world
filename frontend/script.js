const emailForm = document.getElementById("emailForm");
const otpForm = document.getElementById("otpForm");
const messageBox = document.getElementById("message");
let email = "";

function showMessage(msg, type = "success") {
  messageBox.textContent = msg;
  messageBox.className = `message ${type}`;
  messageBox.style.display = "block";
}

emailForm.addEventListener("submit", async (e) => {
  e.preventDefault();
  email = document.getElementById("email").value.trim();

  try {
    const res = await fetch("http://localhost:8080/api/request-otp", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email }),
    });

    const text = await res.text();

    if (res.ok) {
      showMessage("OTP sent to your email.");
      otpForm.style.display = "block";
      emailForm.style.display = "none";
    } else {
      showMessage(text, "error");
    }
  } catch (err) {
    showMessage("Something went wrong. Try again.", "error");
  }
});

otpForm.addEventListener("submit", async (e) => {
  e.preventDefault();
  const otp = document.getElementById("otp").value.trim();

  try {
    const res = await fetch("http://localhost:8080/api/verify-otp", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, otp }),
    });

    const text = await res.text();

    if (res.ok) {
      showMessage("✅ Email verified successfully!");
      otpForm.style.display = "none";
    } else {
      showMessage(text, "error");
    }
  } catch (err) {
    showMessage("Failed to verify OTP.", "error");
  }
});
