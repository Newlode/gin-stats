package stats

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	metrics "github.com/rcrowley/go-metrics"
)

const (
	ginLatencyMetric = "gin.latency"
	ginStatusMetric  = "gin.status"
	ginRequestMetric = "gin.request"
)

//Report from default metric registry
func Report() metrics.Registry {
	return metrics.DefaultRegistry
}

//RequestStats middleware
func RequestStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		req := metrics.GetOrRegisterMeter(ginRequestMetric, nil)
		req.Mark(1)

		latency := metrics.GetOrRegisterTimer(ginLatencyMetric, nil)
		latency.UpdateSince(start)

		handlerName := strings.Replace(c.Request.URL.Path, "/", "-", -1)
		status := metrics.GetOrRegisterMeter(fmt.Sprintf("%s.%s.%d", ginStatusMetric, handlerName, c.Writer.Status()), nil)
		status.Mark(1)
	}
}
