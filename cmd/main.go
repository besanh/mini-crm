package main

import (
	"errors"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
	log "github.com/besanh/logger/logging/slog"
	v1 "github.com/besanh/mini-crm/apis/v1"
	"github.com/besanh/mini-crm/common/cache"
	"github.com/besanh/mini-crm/common/env"
	"github.com/besanh/mini-crm/common/response"
	"github.com/besanh/mini-crm/common/util"
	"github.com/besanh/mini-crm/common/variable"
	"github.com/besanh/mini-crm/pkgs/messagequeue"
	"github.com/besanh/mini-crm/pkgs/mongodb"
	"github.com/besanh/mini-crm/pkgs/oauth"
	"github.com/besanh/mini-crm/pkgs/redis"
	"github.com/besanh/mini-crm/pkgs/sqlclient"
	"github.com/besanh/mini-crm/repositories"
	server "github.com/besanh/mini-crm/servers"
	"github.com/besanh/mini-crm/services"
	"github.com/cardinalby/hureg"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"github.com/joho/godotenv"
)

var (
	DB             mongodb.IMongoDBClient
	sessionManager *scs.SessionManager
)

// init is the entry point of the application
func init() {
	// Load env file
	if err := godotenv.Load("./.env"); err != nil {
		panic(err)
	}

	// Set application version and name
	variable.API_VERSION = env.GetStringENV("API_VERSION", "v1.0")
	variable.API_SERVICE_NAME = env.GetStringENV("API_SERVICE_NAME", "mini_crm")

	// Initialize logs
	initLogger()

	// Initialize redis
	if enableRedis := env.GetBoolENV("ENABLE_REDIS", false); enableRedis {
		initRedis()
	}

	// Initialize MongoDB
	if enableMongo := env.GetBoolENV("ENABLE_MONGODB", false); enableMongo {
		initMongoDb()
	}

	// Initialize NATS JetStream
	if enableNatsJetstream := env.GetBoolENV("ENABLE_NATS", false); enableNatsJetstream {
		initNatsJetstream()
	}

	// Initialize PostgreSQL
	if enablePostgreSql := env.GetBoolENV("ENABLE_PG", false); enablePostgreSql {
		initPostgreSql()
	}
}

// initLogger initializes the logger with the given log level and log file.
// If a log server is provided, it will send the logs to the server.
func initLogger() {
	logFile := "./tmp/console.log"
	logLevel := log.LEVEL_DEBUG
	switch env.GetStringENV("LOG_LEVEL", "error") {
	case "debug":
		logLevel = log.LEVEL_DEBUG
	case "info":
		logLevel = log.LEVEL_INFO
	case "error":
		logLevel = log.LEVEL_ERROR
	case "warn":
		logLevel = log.LEVEL_WARN
	}
	opts := []log.Option{}
	opts = append(opts, log.WithLevel(logLevel),
		log.WithRotateFile(logFile),
		log.WithFileSource(),
		log.WithTraceId(),
		log.WithAttrs(slog.Attr{
			Key: "environment", Value: slog.StringValue(env.GetStringENV("ENVIRONMENT", "local")),
		}),
	)
	// If a log server is provided, send the logs to the server.
	if env.GetStringENV("LOG_SERVER", "") != "" {
		// get server and port from env
		arr := strings.Split(env.GetStringENV("LOG_SERVER", ""), ":")
		if len(arr) >= 2 {
			tag := "fcm"
			client, err := fluent.New(fluent.Config{FluentPort: int(util.ParseInt64(arr[1])), FluentHost: arr[0]})
			if err != nil {
				log.Error(err)
			} else {
				opts = append(opts, log.WithFluentd(client, tag))
			}
		}
	}
	// Set the logger with the given options.
	log.SetLogger(log.NewSLogger(opts...))
}

// initMongoDb initializes the MongoDB client.
// It gets the MongoDB connection string from the environment variables MONGODB_HOST, MONGODB_PORT, MONGODB_DATABASE,
// MONGODB_USERNAME, MONGODB_PASSWORD, and MONGODB_DEFAULT_AUTH_DB.
// If the connection string is invalid, it panics.
func initMongoDb() {
	mongodbConfig := mongodb.MongoDBConfig{
		Username:      env.GetStringENV("MONGODB_USERNAME", ""),
		Password:      env.GetStringENV("MONGODB_PASSWORD", ""),
		Host:          env.GetStringENV("MONGODB_HOST", "localhost"),
		Port:          env.GetIntENV("MONGODB_PORT", 27017),
		Database:      env.GetStringENV("MONGODB_DATABASE", "fcm"),
		DefaultAuthDb: env.GetStringENV("MONGODB_DEFAULT_AUTH_DB", "admin"),
	}

	var err error
	var db mongodb.IMongoDBClient
	db, err = mongodb.NewMongoDBClient(mongodbConfig)
	if err != nil {
		// If the connection string is invalid, panic
		log.Errorf("mongodb connect error: %v", err)
		panic(err)
	}

	DB = db
}

// initRedis initializes the Redis client.
// It gets the Redis connection string from the environment variables REDIS_HOST, REDIS_PASSWORD, REDIS_DB, REDIS_POOL_SIZE,
// REDIS_POOL_TIMEOUT, REDIS_READ_TIMEOUT, and REDIS_WRITE_TIMEOUT.
// If the connection string is invalid, it panics.
func initRedis() {
	redisClient := &redis.RedisConfig{
		Host:         env.GetStringENV("REDIS_HOST", "localhost"),
		Password:     env.GetStringENV("REDIS_PASSWORD", ""),
		DB:           env.GetIntENV("REDIS_DB", 0),
		PoolSize:     env.GetIntENV("REDIS_POOL_SIZE", 10),
		PoolTimeout:  env.GetIntENV("REDIS_POOL_TIMEOUT", 10),
		ReadTimeout:  env.GetIntENV("REDIS_READ_TIMEOUT", 10),
		WriteTimeout: env.GetIntENV("REDIS_WRITE_TIMEOUT", 10),
	}

	var err error
	if redis.Redis, err = redis.NewRedis(*redisClient); err != nil {
		// If the connection string is invalid, panic
		log.Errorf("redis connect error: %v", err)
		panic(err)
	}

	// Initialize the Redis cache with the Redis client.
	cache.RCache = cache.NewRedisCache(redis.Redis.GetClient())

	redisCfg := messagequeue.Rcfg{
		Address:  env.GetStringENV("REDIS_ADDRESS", "localhost"),
		Password: env.GetStringENV("REDIS_PASSWORD", ""),
		DB:       env.GetIntENV("REDIS_RMQ_DB", 1),
	}
	messagequeue.RMQ = messagequeue.NewRMQ(redisCfg)
}

// initNatsJetstream initializes the NATS JetStream client.
// It gets the NATS JetStream connection string from the environment variable NATS_JETSTREAM_HOST.
// If the connection string is invalid, it panics.
func initNatsJetstream() {
	nat := &messagequeue.NatsJetStream{
		Config: messagequeue.Config{
			Host: env.GetStringENV("NATS_JETSTREAM_HOST", "localhost:4222"),
		},
	}

	// Connect to NATS JetStream
	if err := nat.Connect(); err != nil {
		// If the connection string is invalid, panic
		log.Errorf("nats jetstream connect error: %v", err)
		panic(err)
	}
}

func initPostgreSql() {
	sqlClientConfig := sqlclient.SqlConfig{
		Host:         env.GetStringENV("PG_HOST", ""),
		Database:     env.GetStringENV("PG_DB", ""),
		Username:     env.GetStringENV("PG_USERNAME", ""),
		Password:     env.GetStringENV("PG_PASSWORD", ""),
		Port:         env.GetIntENV("PG_PORT", 5432),
		DialTimeout:  20,
		ReadTimeout:  30,
		WriteTimeout: 30,
		Timeout:      30,
		PoolSize:     10,
		MaxOpenConns: 20,
		MaxIdleConns: 10,
		Driver:       sqlclient.POSTGRESQL,
	}
	repositories.DBConn = sqlclient.NewSqlClient(sqlClientConfig)
}

// main is the entry point of the application.
func main() {
	// Decrypt and validate the secret key from environment variables.
	isOk, err := util.DecryptSecret(env.GetStringENV("SECRET_KEY", ""))
	if err != nil {
		panic(err) // Terminate if decryption fails.
	} else if !isOk {
		panic(errors.New("secret_key was incorrect")) // Terminate if the secret key is incorrect.
	}

	// Determine the environment mode for Gin framework.
	// Gin
	envMode := env.GetStringENV("ENV", "debug")
	if slices.Contains([]string{"debug", "test", "release"}, envMode) {
		panic(errors.New("env was incorrect")) // Terminate if the environment mode is invalid.
	}

	// Manage session
	// Initialize session management.
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour // Set session lifetime.
	sessionManager.Cookie.Persist = true     // Make session cookies persistent.
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode

	// Set true if using HTTPS
	// Set to true if using HTTPS to ensure secure cookies.
	sessionManager.Cookie.Secure = false

	router := server.NewServer(envMode, sessionManager)

	// Initialize server
	initServer(router)
	response.NewHumaError()

	server.Start(router, env.GetStringENV("API_PORT", "8000"))
}

// initServer initializes the server with the given Gin engine.
// It creates a new Huma API generator and initializes the services with the server.
func initServer(server *gin.Engine) {
	// Create a new Huma API generator.
	humaAPI := humagin.New(server, huma.Config{
		// OpenAPI configuration.
		OpenAPI: &huma.OpenAPI{
			// OpenAPI version.
			OpenAPI: "3.1.0",
			// API information.
			Info: &huma.Info{
				// API title.
				Title: variable.API_SERVICE_NAME,
				// API version.
				Version: variable.API_VERSION,
				// API contact information.
				Contact: &huma.Contact{
					// Name of the contact person.
					Name: "ANHLE- Mini CRM",
					// URL of the API documentation.
					URL: "https://github.com/besanh/mini-crm",
					// Email address of the contact person.
					Email: "anhle3532@gmail.com",
				},
			},
			// Components of the OpenAPI specification.
			Components: &huma.Components{
				// Security schemes of the API.
				SecuritySchemes: map[string]*huma.SecurityScheme{
					// Bearer authentication.
					"miniCrmAuth": {
						// Type of the security scheme.
						Type: "http",
						// Scheme of the security scheme.
						Scheme: "bearer",
						// Location of the security scheme.
						In: "header",
						// Description of the security scheme.
						Description:  "Authorization header using the Bearer scheme. Example: \"Authorization: Bearer {token}\"",
						BearerFormat: "Token String",
						// Format of the token.
						Name: "Authorization",
						// Name of the security scheme.
					},
				},
			},
		},
		OpenAPIPath: "/anhle/openapi",
		// Path to the OpenAPI specification.
		DocsPath: "",
		// Path to the documentation.
		Formats: huma.DefaultFormats,
		// Supported formats of the API.
		DefaultFormat: "application/json",
		// Default format of the API.
	})

	server.GET("/anhle/openapi/docs", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
		<!doctype html>
		<html>
			<head>
				<title>Mini CRM APIs</title>
				<meta charset="utf-8" />
				<meta name="viewport" content="width=device-width, initial-scale=1" />
			</head>
			<body>
				<script
					id="api-reference"
					data-url="/anhle/openapi.json">
				</script>
				<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
			</body>
		</html>
		`))
	})

	// Create a new Hureg API generator.

	// Initialize the services with the server.
	api := hureg.NewAPIGen(humaAPI)
	initServices(server, &api)
}

// initServices initializes services with the server instance.
// It creates new instances of the repositories, services, and handlers, and registers them with the server.
func initServices(server *gin.Engine, api *hureg.APIGen) {
	// Repositories
	usersRepo := repositories.NewUsers(repositories.DBConn.GetDB())

	// Services
	// Create a new OAuth2 configuration from environment variables.
	oau2Scope := env.GetSliceStringENV("OAUTH2_SCOPE", []string{})
	services.OAUTH2CONFIG = &oauth.OAuth2Config{
		ClientId:     env.GetStringENV("OAUTH2_CLIENT_ID", ""),
		ClientSecret: env.GetStringENV("OAUTH2_CLIENT_SECRET", ""),
		Scopes:       oau2Scope,
		Endpoint: oauth2.Endpoint{
			AuthURL:  env.GetStringENV("OAUTH2_ENDPOINT_AUTH_URL", ""),
			TokenURL: env.GetStringENV("OAUTH2_ENDPOINT_TOKEN_URL", ""),
		},
		Redirect: env.GetStringENV("OAUTH2_REDIRECT_URL", ""),
	}
	// Create a new server instance with the specified environment mode and session manager.
	services.ENABLE_LOGIN_MULTI_SESSION = env.GetBoolENV("ENABLE_LOGIN_MULTI_SESSION", false)

	// Configure service settings from environment variables.
	services.GOOGLE_URL_USER_INFO = env.GetStringENV("GOOGLE_URL_USER_INFO", "")

	// Create a new OAuth2 client from the OAuth2 configuration.
	// oAuth2Client := oauth.NewOAuth2(*services.OAUTH2CONFIG)
	userService := services.NewUsers(usersRepo)

	// Handlers
	v1.NewHealthCheck(api)
	v1.NewUsers(api, userService)
}
