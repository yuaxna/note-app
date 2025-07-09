package backend

import (
	"database/sql"
	"net/http"
	"regexp"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterAuthRoutes(r *gin.Engine) {
	r.POST("/signup", Signup)
	r.POST("/login", Login)
	r.POST("/logout", Logout)
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	return re.MatchString(email)
}

func isValidPassword(password string) bool {
	if len(password) < 5 {
		return false
	}

	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[\W_]`).MatchString(password)

	return hasLower && hasUpper && hasDigit && hasSpecial
}

func Signup(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if user.Fullname == "" || user.Email == "" || user.Username == "" || user.Password == "" || user.Gender == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	if len(user.Username) < 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username must be at least 5 characters"})
		return
	}

	if !isValidEmail(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	if !isValidPassword(user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain upper, lower, number, special char and be 5+ chars"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	_, err = DB.Exec(
		"INSERT INTO users (fullname, email, username, password, gender) VALUES (?, ?, ?, ?, ?)",
		user.Fullname, user.Email, user.Username, string(hashedPassword), user.Gender,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists with that email or username"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Signup successful"})
}

func Login(c *gin.Context) {
	var login User
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	row := DB.QueryRow(
		"SELECT id, password, fullname, email, username, gender FROM users WHERE username = ? OR email = ?",
		login.Identifier, login.Identifier,
	)

	var storedID int
	var storedHashedPassword, fullname, email, username, gender string
	err := row.Scan(&storedID, &storedHashedPassword, &fullname, &email, &username, &gender)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(login.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Set session - now including gender
	session := sessions.Default(c)
	session.Set("user_id", storedID)
	session.Set("fullname", fullname)
	session.Set("email", email)
	session.Set("username", username)
	session.Set("gender", gender) // Add this line
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func GetMe(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Fetch fresh user data from database
	row := DB.QueryRow(
		"SELECT id, fullname, email, username, COALESCE(gender, 'Not specified') as gender FROM users WHERE id = ?",
		userID,
	)

	var id int
	var fullname, email, username, gender string
	err := row.Scan(&id, &fullname, &email, &username, &gender)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       id,
		"fullname": fullname,
		"email":    email,
		"username": username,
		"gender":   gender,
	})
}

func isAPIRequest(c *gin.Context) bool {
	return c.Request.Header.Get("Content-Type") == "application/json" ||
		c.Request.Header.Get("Accept") == "application/json" ||
		(len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api") ||
		c.Request.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

// Temporary debug function - remove after fixing
func DebugUser(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check what's in the database
	row := DB.QueryRow("SELECT id, fullname, email, username, gender FROM users WHERE id = ?", userID)

	var id int
	var fullname, email, username, gender sql.NullString
	err := row.Scan(&id, &fullname, &email, &username, &gender)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"database_data": gin.H{
			"id":           id,
			"fullname":     fullname.String,
			"email":        email.String,
			"username":     username.String,
			"gender":       gender.String,
			"gender_valid": gender.Valid,
		},
		"session_data": gin.H{
			"user_id":  session.Get("user_id"),
			"fullname": session.Get("fullname"),
			"email":    session.Get("email"),
			"username": session.Get("username"),
			"gender":   session.Get("gender"),
		},
	})
}
