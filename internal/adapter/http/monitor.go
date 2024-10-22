package http

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type MonitoringHandler struct {
	MongoClient *mongo.Client
}

func NewMonitoringHandler(client *mongo.Client) *MonitoringHandler {
	return &MonitoringHandler{MongoClient: client}
}

func (h *MonitoringHandler) GetMonitoringData(c *gin.Context) {
	collection := h.MongoClient.Database("Monitoring").Collection("logs")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, map[string]interface{}{})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var results []map[string]interface{}
	if err = cursor.All(ctx, &results); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, results)
}
