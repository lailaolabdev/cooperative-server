package main

import (
	"context"
	"cooperative-service/internal/config"
	"cooperative-service/internal/database"
	"cooperative-service/internal/middleware"
	authmodule "cooperative-service/internal/modules/auth"
	coopmodule "cooperative-service/internal/modules/cooperative"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	client, err := database.Connect(cfg.MongoURI)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	repo := coopmodule.NewRepository(client.Database(cfg.Database))
	authRepo := authmodule.NewRepository(client.Database(cfg.Database))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = repo.EnsureIndexes(ctx); err != nil {
		log.Fatal(err)
	}
	if err = authRepo.EnsureIndexes(ctx); err != nil {
		log.Fatal(err)
	}
	authService := authmodule.NewService(authRepo, cfg.JWTSecret)
	authHandler := authmodule.NewHandler(authService)
	coopHandler := coopmodule.NewHandler(coopmodule.NewService(repo))
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), cors.New(cors.Config{AllowOrigins: cfg.AllowedOrigins, AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, AllowHeaders: []string{"Origin", "Content-Type", "Authorization"}, MaxAge: 12 * time.Hour}))
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	api := r.Group("/api/v1")
	api.POST("/auth/admin/login", authHandler.Login)
	api.GET("/cooperatives", coopHandler.List)
	api.GET("/cooperatives/:id", coopHandler.Get)
	admin := api.Group("/admin", middleware.AdminAuth(cfg.JWTSecret))
	admin.POST("/cooperatives", coopHandler.Create)
	admin.PUT("/cooperatives/:id", coopHandler.Update)
	admin.DELETE("/cooperatives/:id", coopHandler.Delete)
	log.Printf("API listening on :%s", cfg.Port)
	if err = r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
