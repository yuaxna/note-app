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
	// Initialize database
	backend.InitDB()

	// Setup Gin router
	router := gin.Default()

	// Setup session store (make sure to change the secret in production)
	store := cookie.NewStore([]byte("your-very-secret-key")) // TODO: Replace with env/config
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})
	router.Use(sessions.Sessions("olive_session", store))

	// Serve static files
	router.Static("/css", "./frontend/css")
	router.Static("/js", "./frontend/js")
	router.LoadHTMLGlob("frontend/templates/*")

	// Auth (public)
	backend.RegisterAuthRoutes(router)

	// Public HTML Routes
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "auths.html", nil)
	})

	// Protected HTML Pages
	securedPages := router.Group("/")
	securedPages.Use(backend.AuthRequired())
	{
		securedPages.GET("/home", func(c *gin.Context) {
			c.HTML(http.StatusOK, "home.html", nil)
		})
		securedPages.GET("/profile", func(c *gin.Context) {
			c.HTML(http.StatusOK, "profile.html", nil)
		})
		securedPages.GET("/shared", func(c *gin.Context) {
			c.HTML(http.StatusOK, "shared.html", nil)
		})

	}

	// Protected API Routes
	api := router.Group("/api")
	api.Use(backend.AuthRequired())
	{
		api.GET("/me", backend.GetMe)
		api.GET("/notes", backend.GetUserNotes)
		api.GET("/debug", backend.DebugUser)
		api.POST("/notes", backend.CreateNote)
		api.PUT("/notes", backend.UpdateNote)
		api.DELETE("/notes/:id", backend.DeleteNote)
		api.POST("/share", backend.ShareNote)
		api.GET("/shared", backend.GetSharedNotes)
		api.POST("/notes/share", backend.ShareNote)
		api.GET("/users", backend.AuthRequired(), backend.GetUsers)
	}

	// WebSocket endpoint for real-time collaboration (you will implement handler)
	router.GET("/ws", backend.WSHandler)

	// Start Server
	log.Println("Server running at http://localhost:8080")
	router.Run(":8080")
}
