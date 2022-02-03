package application

import (
	"github.com/gin-gonic/gin"
	"github.com/zzz4zzz/go-playground/db"
	"github.com/zzz4zzz/go-playground/http/controllers"
	"github.com/zzz4zzz/go-playground/http/middlewares"
	"github.com/zzz4zzz/go-playground/trace"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"gorm.io/gorm"
)

type Application struct {
	DatabaseConfiguration db.DatabaseConfiguration
	GinEngine             *gin.Engine
	Router                gin.IRouter
	DB                    *gorm.DB
	Address               string
	TracerProvider        *trace.Provider
}

func (a *Application) Run() error {
	return a.GinEngine.Run("localhost:8080")
}

func (a *Application) configureRoutes() {
	a.Router.Use(
		middlewares.TransactionMiddleware(a.DB),
		otelgin.Middleware("playground", otelgin.WithTracerProvider(a.TracerProvider.Provider), otelgin.WithPropagators(a.TracerProvider.Propagators)),
	)

	a.Router.POST("/users", controllers.AddUserTraced(a.TracerProvider))

	graphqlGroup := a.Router.Group("/graphql")
	graphqlGroup.Use()
	{
		graphqlGroup.POST("", controllers.GinGraphQLHandler)
		graphqlGroup.GET("", controllers.GinGraphQLHandler)
	}
}
func (a *Application) connectToDatabase() {
	a.DB = db.InitDatabase(a.DatabaseConfiguration)
}
func (a *Application) configureTraceProvider() {
	provider, err := trace.NewProvider(trace.JaegerConfiguration{
		Endpoint:    "http://localhost:14268/api/traces",
		ServiceName: "playground",
		Enabled:     true,
	})
	if err != nil {
		panic("Trace provider not configured")
	}
	a.TracerProvider = &provider
}

func NewApplication(databaseConfiguration db.DatabaseConfiguration, address string) *Application {
	r := gin.Default()
	application := &Application{
		DatabaseConfiguration: databaseConfiguration,
		GinEngine:             r,
		Router:                r,
		Address:               address,
	}
	application.configureTraceProvider()
	application.connectToDatabase()
	application.configureRoutes()
	return application
}
