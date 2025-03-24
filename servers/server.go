package server

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/besanh/mini-crm/common/log"
	"github.com/besanh/mini-crm/common/variable"
	"github.com/gin-gonic/gin"
)

func NewServer(envMode string, sessionManager *scs.SessionManager) *gin.Engine {
	switch envMode {
	case "test":
		gin.SetMode(gin.TestMode)
	case "release":
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()

	// Use the session middleware
	engine.Use(func(c *gin.Context) {
		sessionManager.LoadAndSave(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).
			ServeHTTP(c.Writer, c.Request)
		c.Next()
	})

	engine.Use(gin.Recovery())
	engine.MaxMultipartMemory = 100 << 20
	engine.Use(CORSMiddleware())
	engine.Use(allowOptionsMethod())
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": variable.API_SERVICE_NAME,
			"version": variable.API_VERSION,
			"time":    time.Now().Unix(),
		})
	})
	return engine
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
		c.Next()
	}
}

func allowOptionsMethod() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}

func Start(server *gin.Engine, port string) {
	v := make(chan struct{})
	go func() {
		if err := server.Run(":" + port); err != nil {
			log.Error("failed to start service")
			close(v)
		}
	}()
	log.Infof("service %v listening on port %v", variable.API_SERVICE_NAME, port)
	<-v
}
