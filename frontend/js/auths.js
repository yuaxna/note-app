document.addEventListener('DOMContentLoaded', () => {
    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form');
    const showRegisterBtn = document.getElementById('show-register');
    const showLoginBtn = document.getElementById('show-login');

    const loginErrorDiv = document.getElementById('login-error');
    const registerErrorDiv = document.getElementById('register-error');

    const logoutBtn = document.getElementById("logout-btn");
    if (logoutBtn) {
        logoutBtn.addEventListener("click", async () => {
            try {
                const res = await fetch("http://localhost:8080/logout", {
                    method: "POST",
                    credentials: "include"
                });

                const data = await res.json();

                if (res.ok) {
                    window.location.href = "/";
                } else {
                    console.error(data.error || "Logout failed.");
                }
            } catch (err) {
                console.error("Error during logout:", err);
            }
        });
    }

    // Helper validation functions
    function validateEmail(email) {
        const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return re.test(email);
    }

    function validatePassword(password) {
        // Min 5 chars, at least 1 uppercase, 1 lowercase, 1 number, 1 special char
        const re = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[\W_]).{5,}$/;
        return re.test(password);
    }

    // Toggle forms
    showRegisterBtn.addEventListener('click', () => {
        loginErrorDiv.textContent = "";
        registerErrorDiv.textContent = "";
        loginForm.style.display = 'none';
        registerForm.style.display = 'block';
    });

    showLoginBtn.addEventListener('click', () => {
        loginErrorDiv.textContent = "";
        registerErrorDiv.textContent = "";
        registerForm.style.display = 'none';
        loginForm.style.display = 'block';
    });

    // Handle login submit
    loginForm.addEventListener("submit", async function (e) {
        e.preventDefault();

        const identifier = document.getElementById("login-identifier").value.trim();
        const password = document.getElementById("login-password").value;
        loginErrorDiv.textContent = "";

        if (!identifier || !password) {
            loginErrorDiv.textContent = "Please enter both your email/username and password.";
            return;
        }

        const res = await fetch("http://localhost:8080/login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify({ identifier, password })
        });

        const data = await res.json();

        if (res.ok) {
            window.location.href = "/home";
        } else {
            loginErrorDiv.textContent = data.error || "Login failed.";
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

        registerErrorDiv.textContent = "";

        if (!fullname || !email || !username || !password || !passwordConfirm || !gender) {
            registerErrorDiv.textContent = 'Please fill all the fields.';
            return;
        }

        if (username.length < 5) {
            registerErrorDiv.textContent = 'Username must be at least 5 characters.';
            return;
        }

        if (!validateEmail(email)) {
            registerErrorDiv.textContent = 'Invalid email format.';
            return;
        }

        if (!validatePassword(password)) {
            registerErrorDiv.textContent = 'Password must be at least 5 characters and include uppercase, lowercase, number, and special character.';
            return;
        }

        if (password !== passwordConfirm) {
            registerErrorDiv.textContent = 'Passwords do not match!';
            return;
        }

        const res = await fetch('http://localhost:8080/signup', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: "include",
            body: JSON.stringify({ fullname, email, username, password, gender }),
        });

        const data = await res.json();

        if (res.ok) {
            registerErrorDiv.style.color = "green";
            registerErrorDiv.textContent = "Registration successful. You can now log in.";
            setTimeout(() => {
                registerErrorDiv.textContent = "";
                showLoginBtn.click(); // switch back to login form
            }, 1000);
        } else {
            registerErrorDiv.textContent = data.error || 'Registration failed.';
        }
    });
});
