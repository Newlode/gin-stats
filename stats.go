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
		handlerName := strings.Replace(c.Request.URL.Path, "/", "-", -1)
		handlerName = handlerName[1:len(handlerName)]

		// Requests Per Second (total and per-handler)
		totalReq := metrics.GetOrRegisterMeter(ginRequestMetric, nil)
		totalReq.Mark(1)

		req := metrics.GetOrRegisterMeter(fmt.Sprintf("%s.%s", ginRequestMetric, handlerName), nil)
		req.Mark(1)
		
		// Latency (total and per-handler)
		totalLatency := metrics.GetOrRegisterTimer(ginLatencyMetric, nil)
		totalLatency.UpdateSince(start)

		latency := metrics.GetOrRegisterTimer(fmt.Sprintf("%s.%s", ginLatencyMetric, handlerName), nil)
		latency.UpdateSince(start)

		// HTTP Status, e.g. 200, 404, 500 (total and per-handler)
		totalStatus := metrics.GetOrRegisterMeter(fmt.Sprintf("%s.%d", ginStatusMetric, handlerName, c.Writer.Status()), nil)
		totalStatus.Mark(1)

		status := metrics.GetOrRegisterMeter(fmt.Sprintf("%s.%s.%d", ginStatusMetric, handlerName, c.Writer.Status()), nil)
		status.Mark(1)
	}
}
