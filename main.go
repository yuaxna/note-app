package main

import (
	"log"
	"net/http"
	"note-app/backend"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database
	backend.InitDB()

	// Create router
	router := gin.Default()

	store := cookie.NewStore([]byte("your-secret-key"))
	router.Use(sessions.Sessions("olive_session", store))

	// Serve static files
	router.Static("/css", "./frontend/css")
	router.Static("/js", "./frontend/js")

	// Load templates
	router.LoadHTMLGlob("frontend/templates/*")

	// Register backend auth routes (/signup, /login)
	backend.RegisterAuthRoutes(router)

	// Route: root → login/register page
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "auths.html", nil)
	})

	// Route: home page
	router.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", nil)
	})

	// Route: profile page
	router.GET("/profile", func(c *gin.Context) {
		c.HTML(http.StatusOK, "profile.html", nil)
	})

	// API route: create a note
	router.POST("/api/notes", backend.CreateNoteHandler)
	router.GET("/api/notes", backend.GetUserNotes)
	router.GET("/api/me", backend.AuthRequired(), backend.GetMe)
	router.POST("/logout", backend.Logout)

	log.Println("Server running at http://localhost:8080")
	router.Run(":8080")
}
