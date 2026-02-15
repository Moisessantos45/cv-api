package main

import (
	"context"
	"cv_api/config"
	"cv_api/internal/middleware"
	"cv_api/internal/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

var ctx = context.Background()

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println("Environment variables loaded successfully")

	API_URL := os.Getenv("SUPABASE_URL")
	API_KEY := os.Getenv("SUPABASE_KEY")
	HOST_URL_DEV := os.Getenv("HOST_URL_DEV")
	HOST_URL_PROD := os.Getenv("HOST_URL_PROD")
	HOST_URL_PROD_WWW := os.Getenv("HOST_URL_PROD_WWW")

	log.Printf("HOST_URL_DEV: %s", HOST_URL_DEV)
	log.Printf("HOST_URL_PROD: %s", HOST_URL_PROD)
	log.Printf("HOST_URL_PROD_WWW: %s", HOST_URL_PROD_WWW)

	err = config.Init(API_URL, API_KEY)
	if err != nil {
		log.Fatal("Error initializing Supabase client: ", err)
	}

	err = config.InitRedis(ctx)
	if err != nil {
		log.Fatal("Error initializing Redis client: ", err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{HOST_URL_DEV, HOST_URL_PROD, HOST_URL_PROD_WWW},
		AllowMethods:     []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.Use(middleware.RateLimiterMiddleware(rate.Every(time.Minute/10), 10))

	api := r.Group("/api/v1")
	{
		routes.ProjectRoutes(api)
		routes.VideoRoutes(api)
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the CV API",
		})
	})

	r.GET("/api/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK", "ip": c.ClientIP()})
	})

	var wg sync.WaitGroup
	wg.Go(func() {
		middleware.StartCleanup()
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

	wg.Wait()
	log.Println("Server exiting")
}
