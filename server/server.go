package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/Depado/ginprom"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/graph-gophers/graphql-go"

	"github.com/lukaszbudnik/migrator/common"
	"github.com/lukaszbudnik/migrator/config"
	"github.com/lukaszbudnik/migrator/converter"
	"github.com/lukaszbudnik/migrator/coordinator"
	"github.com/lukaszbudnik/migrator/data"
	"github.com/lukaszbudnik/migrator/metrics"
	"github.com/lukaszbudnik/migrator/types"
)

const (
	defaultPort     string = "8080"
	requestIDHeader string = "X-Request-ID"
)

type errorResponse struct {
	Errors []errorMessage `json:"errors"`
}

type errorMessage struct {
	Message string `json:"message"`
}

// GetPort gets the port from config or defaultPort
func GetPort(config *config.Config) string {
	if strings.TrimSpace(config.Port) == "" {
		return defaultPort
	}
	return config.Port
}

func requestIDHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get(requestIDHeader)
		if requestID == "" {
			requestID = fmt.Sprintf("%d", time.Now().UnixNano())
		}
		ctx := context.WithValue(c.Request.Context(), common.RequestIDKey{}, requestID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func logLevelHandler(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), common.LogLevelKey{}, config.LogLevel)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				common.LogPanic(c.Request.Context(), "Panic recovered: %v", err)
				if gin.IsDebugging() {
					debug.PrintStack()
				}
				errorMsg := errorMessage{err.(string)}
				c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{Errors: []errorMessage{errorMsg}})
			}
		}()
		c.Next()
	}
}

func requestLoggerHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		common.LogInfo(c.Request.Context(), "clientIP=%v method=%v request=%v", c.ClientIP(), c.Request.Method, c.Request.URL.RequestURI())
		c.Next()
	}
}

func deprecationHeaderHandler(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if deprecated config fields are being used
		if config.IsUsingDeprecatedTenantSelectSQL() || config.IsUsingDeprecatedTenantInsertSQL() {
			c.Header("Deprecation", "true")
			c.Header("Sunset", "Thu, 31 Dec 2026 23:59:59 GMT")
			c.Header("Warning", "299 - \"tenantSelectSQL and tenantInsertSQL are deprecated since v2025.1.0. Use tenantSelect and tenantInsert instead. These fields will be removed in v2027.0.0.\"")
		}
		c.Next()
	}
}

func makeHandler(config *config.Config, metrics metrics.Metrics, newCoordinator coordinator.Factory, handler func(*gin.Context, *config.Config, metrics.Metrics, coordinator.Factory)) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(c, config, metrics, newCoordinator)
	}
}

func configHandler(c *gin.Context, config *config.Config, metrics metrics.Metrics, newCoordinator coordinator.Factory) {
	c.YAML(200, config)
}

func schemaHandler(c *gin.Context, config *config.Config, metrics metrics.Metrics, newCoordinator coordinator.Factory) {
	c.String(http.StatusOK, strings.TrimSpace(data.SchemaDefinition))
}

func healthHandler(c *gin.Context, config *config.Config, metrics metrics.Metrics, newCoordinator coordinator.Factory) {
	coordinator := newCoordinator(c.Request.Context(), config, metrics)
	healthStatus := coordinator.HealthCheck()

	status := http.StatusOK
	if healthStatus.Status == types.HealthStatusDown {
		status = http.StatusServiceUnavailable
	}

	c.JSON(status, healthStatus)
}

// GraphQL endpoint
func serviceHandler(c *gin.Context, config *config.Config, metrics metrics.Metrics, newCoordinator coordinator.Factory) {
	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		common.LogError(c.Request.Context(), "Bad request: %v", err.Error())
		errorMsg := errorMessage{"Invalid request, please see documentation for valid JSON payload"}
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{Errors: []errorMessage{errorMsg}})
		return
	}

	coordinator := newCoordinator(c.Request.Context(), config, metrics)
	defer coordinator.Dispose()
	opts := []graphql.SchemaOpt{graphql.UseFieldResolvers()}
	schema := graphql.MustParseSchema(data.SchemaDefinition, &data.RootResolver{Coordinator: coordinator}, opts...)

	response := schema.Exec(c.Request.Context(), params.Query, params.OperationName, params.Variables)
	if response.Errors == nil {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusInternalServerError, response)
	}

}

// Handle file upload
func uploadHandler(c *gin.Context, config *config.Config, metrics metrics.Metrics, newCoordinator coordinator.Factory) {
	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		common.LogError(c.Request.Context(), "Failed to get uploaded file: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload file: " + err.Error()})
		return
	}

	// Get conversion parameters
	sourceFormat := c.PostForm("sourceFormat")
	targetFormat := c.PostForm("targetFormat")

	// Create uploads directory if it doesn't exist
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		os.Mkdir("./uploads", 0755)
	}

	// Save the original file
	filename := fmt.Sprintf("./uploads/%s", file.Filename)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		common.LogError(c.Request.Context(), "Failed to save file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file: " + err.Error()})
		return
	}

	// Process the file using our converter
	result := fmt.Sprintf("File %s converted from %s to %s", file.Filename, sourceFormat, targetFormat)
	
	// Convert the file
	convertedData, err := converter.ConvertFile(filename, targetFormat)
	if err != nil {
		common.LogError(c.Request.Context(), "Failed to convert file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert file: " + err.Error()})
		return
	}

	// Save converted data to a file with a unique name
	timestamp := time.Now().Unix()
	migratedFilename := fmt.Sprintf("./uploads/migrated_%d_%s.%s", timestamp, file.Filename, targetFormat)
	if targetFormat == "excel" {
		migratedFilename = fmt.Sprintf("./uploads/migrated_%d_%s.xlsx", timestamp, file.Filename)
	}
	
	err = os.WriteFile(migratedFilename, convertedData, 0644)
	if err != nil {
		common.LogError(c.Request.Context(), "Failed to save converted file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save converted file: " + err.Error()})
		return
	}

	common.LogInfo(c.Request.Context(), "File uploaded and processed: %s", filename)

	c.JSON(http.StatusOK, gin.H{
		"message": result,
		"file":    file.Filename,
		"source":  sourceFormat,
		"target":  targetFormat,
		"migratedFile": fmt.Sprintf("migrated_%d_%s.%s", timestamp, file.Filename, targetFormat),
	})
}

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Handle migrated file download
func downloadMigratedFileHandler(c *gin.Context, config *config.Config, metrics metrics.Metrics, newCoordinator coordinator.Factory) {
	filename := c.Query("filename")
	
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing filename parameter"})
		return
	}
	
	// Construct the full path to the migrated file
	migratedFilePath := fmt.Sprintf("./uploads/%s", filename)
	
	// Check if file exists
	if _, err := os.Stat(migratedFilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Migrated file not found: %s", filename)})
		return
	}
	
	// Serve the file
	c.File(migratedFilePath)
}

func CreateRouterAndPrometheus(versionInfo *types.VersionInfo, config *config.Config, newCoordinator coordinator.Factory) *gin.Engine {
	r := gin.New()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	p := ginprom.New(
		ginprom.Engine(r),
		ginprom.Namespace("migrator"),
		ginprom.Subsystem("gin"),
		ginprom.Path("/metrics"),
	)
	p.AddCustomGauge("info", "Information about migrator app", []string{"version"})
	p.AddCustomGauge("versions_created", "Number of versions created by migrator", []string{})
	p.AddCustomGauge("tenants_created", "Number of migrations applied by migrator", []string{})
	p.AddCustomGauge("migrations_applied", "Number of migrations applied by migrator", []string{"type"})

	p.SetGaugeValue("info", []string{versionInfo.Release + " @ " + versionInfo.Sha}, 1)

	r.Use(p.Instrument())

	metrics := metrics.New(p)

	return SetupRouter(r, versionInfo, config, metrics, newCoordinator)
}

// SetupRouter setups router
func SetupRouter(r *gin.Engine, versionInfo *types.VersionInfo, config *config.Config, metrics metrics.Metrics, newCoordinator coordinator.Factory) *gin.Engine {
	r.HandleMethodNotAllowed = true
	r.Use(logLevelHandler(config), recovery(), requestIDHandler(), requestLoggerHandler(), deprecationHeaderHandler(config))

	// Serve static files
	r.Static("/static", "./static")

	if strings.TrimSpace(config.PathPrefix) == "" {
		config.PathPrefix = "/"
	}

	r.GET(config.PathPrefix+"/", func(c *gin.Context) {
		// Serve the main index.html file
		c.File("./static/index.html")
	})

	r.GET(config.PathPrefix+"/health", makeHandler(config, metrics, newCoordinator, healthHandler))

	// Handle file uploads
	r.POST(config.PathPrefix+"/upload", makeHandler(config, metrics, newCoordinator, uploadHandler))
	
	// Handle migrated file downloads
	r.GET(config.PathPrefix+"/download-migrated", makeHandler(config, metrics, newCoordinator, downloadMigratedFileHandler))

	v1 := r.Group(config.PathPrefix + "/v1")
	v1.Any("/*any", func(c *gin.Context) {
		c.Status(http.StatusGone)
	})

	v2 := r.Group(config.PathPrefix + "/v2")
	v2.GET("/config", makeHandler(config, metrics, newCoordinator, configHandler))
	v2.GET("/schema", makeHandler(config, metrics, newCoordinator, schemaHandler))
	v2.POST("/service", makeHandler(config, metrics, newCoordinator, serviceHandler))

	return r
}