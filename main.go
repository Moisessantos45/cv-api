package main

import (
	"context"
	"cv_api/config"
	"cv_api/config/db"
	"cv_api/internal/routes"
	"cv_api/internal/shared/middleware"
	"cv_api/internal/shared/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielkov/gin-helmet/ginhelmet"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Advertencia: no se encontró .env, se usan variables de entorno del sistema")
	}

	if err := db.Connect(); err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}

	if err := db.InitializeDatabase(); err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}

	if err := config.InitRedis(context.Background()); err != nil {
		fmt.Println("Error initializing Redis:", err)
		return
	}

	ip := utils.GetOutboundIP()

	log.Println("Environment variables loaded successfully", ip)

	HOST_URL_DEV := os.Getenv("HOST_URL_DEV")
	HOST_URL_PROD := os.Getenv("HOST_URL_PROD")
	HOST_URL_PROD_WWW := os.Getenv("HOST_URL_PROD_WWW")
	HOST_API_PROD := os.Getenv("HOST_API_PROD")
	HOST_API_PROD_WWW := os.Getenv("HOST_API_PROD_WWW")

	log.Printf("HOST_URL_DEV: %s", HOST_URL_DEV)
	log.Printf("HOST_URL_PROD: %s", HOST_URL_PROD)
	log.Printf("HOST_URL_PROD_WWW: %s", HOST_URL_PROD_WWW)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(ginhelmet.Default())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{HOST_URL_DEV, HOST_URL_PROD, HOST_URL_PROD_WWW, HOST_API_PROD, HOST_API_PROD_WWW},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.Use(middleware.RateLimiterMiddleware(rate.Every(time.Minute/10), 10))

	routes.Init()
	api := r.Group("/api/v1")
	{
		routes.ProfileRoutes(api)
		routes.ProjectRoutes(api)
		routes.VideoRoutes(api)
	}

	auth := r.Group("/api/v1/auth")
	{
		routes.AuthRoutes(auth)
	}

	routes.PostRoutes(api)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the CV API",
		})
	})

	r.GET("/api/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK", "ip": c.ClientIP()})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Println("Server starting on :4100...")
	srv := &http.Server{
		Addr:    ":4100",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

// CREATE INDEX idx_projects_created_at ON projects(created_at DESC);
// CREATE INDEX idx_projects_state_id ON projects(state_id);
