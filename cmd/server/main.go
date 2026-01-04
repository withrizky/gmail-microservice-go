package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"gmail_microservice/internal/model"
	"gmail_microservice/internal/worker"
)

func main() {
	godotenv.Load()

	// Ambil konfigurasi akun (Format JSON Array String)
	accountsConfig := os.Getenv("GMAIL_ACCOUNTS")

	// Inisialisasi Dispatcher dengan Config Akun
	dispatcher := worker.NewDispatcher(20, 5000, accountsConfig)
	dispatcher.Run()

	r := gin.New()
	r.Use(gin.Recovery())

	r.POST("/send-email", func(c *gin.Context) {
		var req model.EmailPayload
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		select {
		case dispatcher.JobQueue <- req:
			c.JSON(202, gin.H{"status": "queued", "message": "Email dalam antrean"})
		default:
			c.JSON(503, gin.H{"error": "Server sibuk"})
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	srv := &http.Server{Addr: ":" + port, Handler: r}

	go func() {
		log.Printf("Gmail Rotator Service running on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	dispatcher.Stop()
}
