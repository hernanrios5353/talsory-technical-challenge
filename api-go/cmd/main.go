package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/talsory/coding-challenge/api-go/internal/client"
	"github.com/talsory/coding-challenge/api-go/internal/handler"
	"github.com/talsory/coding-challenge/api-go/internal/service"
)

func main() {
	port := getEnv("PORT", "8080")
	nodeAPIURL := getEnv("NODE_API_URL", "http://localhost:3000")

	// Inicializar dependencias
	nodeClient := client.NewNodeClient(nodeAPIURL)
	qrService := service.NewQRService(nodeClient)
	qrHandler := handler.NewQRHandler(qrService)

	authService := service.NewAuthService()
	authHandler := handler.NewAuthHandler(authService)

	app := fiber.New(fiber.Config{
		AppName: "Talsory Technical challenge - API Go",
	})
 
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"service": "api-go",
		})
	})

	// Rutas públicas
	api := app.Group("/api/v1")
	api.Post("/auth/login", authHandler.Login)

	// Rutas protegidas por JWT
	protected := api.Group("", handler.JWTMiddleware)
	protected.Post("/qr", qrHandler.Handle)

	log.Printf("🚀 API Go escuchando en el puerto %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}