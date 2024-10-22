package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-gin-hex-arch/internal/adapter/config"
	"go-gin-hex-arch/internal/util"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
	"time"
)

func MonitoringFuncPerformance(client *mongo.Client, cfg *config.Container) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == cfg.HTTP.Prefix+"/monitoring" {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()
		duration := time.Since(start)

		mongoDB := client.Database(cfg.MONGO.MongoDB).Collection(cfg.MONGO.MongoCollection)

		go func() {
			result, err := mongoDB.InsertOne(context.Background(), map[string]interface{}{
				"path":     c.Request.URL.Path,
				"method":   c.Request.Method,
				"duration": util.FormatDuration(duration),
				"time":     time.Now().In(util.Loc),
			})
			if err != nil {
				slog.Error("Failed to insert performance data", "error", err)
			} else {
				slog.Info("Inserted performance data", "result", result)
			}
		}()
	}
}
