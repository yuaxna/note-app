document.addEventListener('DOMContentLoaded', () => {
    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form');
    const showRegisterBtn = document.getElementById('show-register');
    const showLoginBtn = document.getElementById('show-login');

    // Toggle forms
    showRegisterBtn.addEventListener('click', () => {
        loginForm.style.display = 'none';
        registerForm.style.display = 'block';
    });

    showLoginBtn.addEventListener('click', () => {
        registerForm.style.display = 'none';
        loginForm.style.display = 'block';
    });

    // Handle login submit
    loginForm.addEventListener("submit", async function (e) {
        e.preventDefault();

        const identifier = document.getElementById("login-identifier").value.trim();
        const password = document.getElementById("login-password").value;

        const res = await fetch("http://localhost:8080/login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ identifier, password })
        });

        const data = await res.json();

        if (res.ok) {
            alert("Login successful!");
            window.location.href = "/home"; // Redirect here
        } else {
            alert(data.error || "Login failed.");
        }
    });

    // Handle register submit
    registerForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const fullname = document.getElementById('register-fullname').value.trim();
        const email = document.getElementById('register-email').value.trim();
        const username = document.getElementById('register-username').value.trim();
        const password = document.getElementById('register-password').value;
        const passwordConfirm = document.getElementById('register-password-confirm').value;
        const gender = document.getElementById('register-gender').value;

        if (password !== passwordConfirm) {
            alert('Passwords do not match!');
            return;
        }

        const res = await fetch('http://localhost:8080/signup', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ fullname, email, username, password, gender }),
        });

        const data = await res.json();

        if (res.ok) {
            alert('Registration successful. You can now log in.');
            showLoginBtn.click(); // switch back to login form
        } else {
            alert(data.error || 'Registration failed.');
        }
    });
});
