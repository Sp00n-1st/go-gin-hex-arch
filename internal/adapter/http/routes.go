package http

import (
	"github.com/gin-gonic/gin"
	"go-gin-hex-arch/internal/adapter/config"
	"go-gin-hex-arch/internal/core/service"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(app *gin.Engine, productService *service.ProductService, mongoClient *mongo.Client, cfg *config.HTTP) {
	productHandler := NewProductHandler(*productService)
	monitoringHandler := NewMonitoringHandler(mongoClient)

	app.POST(cfg.Prefix+"/products", productHandler.CreateProduct)
	app.PUT(cfg.Prefix+"/products/:id", productHandler.UpdateProduct)
	app.DELETE(cfg.Prefix+"/products/:id", productHandler.DeleteProduct)
	app.GET(cfg.Prefix+"/products/:id", productHandler.GetProductByID)
	app.GET(cfg.Prefix+"/products", productHandler.GetProducts)

	app.GET(cfg.Prefix+"/monitoring", monitoringHandler.GetMonitoringData)
}
