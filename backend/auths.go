package backend

import (
	"database/sql"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterAuthRoutes(r *gin.Engine) {
	r.POST("/signup", Signup)
	r.POST("/login", Login)
	// r.GET("/api/me", GetMe)
}

func Signup(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email or username already exists"})
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
		"SELECT id, password, fullname, email, username FROM users WHERE username = ? OR email = ?",
		login.Identifier, login.Identifier,
	)

	var storedID int
	var storedHashedPassword, fullname, email, username string
	err := row.Scan(&storedID, &storedHashedPassword, &fullname, &email, &username)
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

	// Set session
	session := sessions.Default(c)
	session.Set("user_id", storedID)
	session.Set("fullname", fullname)
	session.Set("email", email)
	session.Set("username", username)
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// GET /api/me
func GetMe(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// You can fetch more user info from DB if needed
	// or store it in session like we did in login
	fullname := session.Get("fullname")
	email := session.Get("email")
	username := session.Get("username")

	c.JSON(http.StatusOK, gin.H{
		"id":       userID,
		"fullname": fullname,
		"email":    email,
		"username": username,
	})
}

func GetUserNotes(c *gin.Context) {
	userID := c.Query("user_id")

	rows, err := DB.Query("SELECT id, title, content FROM notes WHERE user_id = ?", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
		return
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content); err != nil {
			continue
		}
		notes = append(notes, note)
	}

	c.JSON(http.StatusOK, notes)
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Set("user_id", userID)
		c.Next()
	}
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}
