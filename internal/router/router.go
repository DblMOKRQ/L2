package router

import (
	"awesomeProject/internal/router/handlers"
	"awesomeProject/internal/router/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Router struct {
	rout    *gin.Engine
	handler *handlers.CalendarHandler
	log     *zap.Logger
}

func NewRouter(handler *handlers.CalendarHandler, mode string, log *zap.Logger) *Router {
	switch mode {
	case "debug":
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}
	router := &Router{
		rout:    gin.Default(),
		handler: handler,
		log:     log,
	}
	router.setupRouter()

	return router
}

func (r *Router) setupRouter() {
	r.rout.Use(middleware.LoggingMiddleware(r.log))
	r.rout.POST("/create_event", r.handler.CreateEvent)
	r.rout.POST("/update_event", r.handler.UpdateEvent)
	r.rout.POST("/delete_event", r.handler.DeleteEvent)
	r.rout.GET("/events_for_day", r.handler.GetEventsForDay)
	r.rout.GET("/events_for_week", r.handler.GetEventsForWeek)
	r.rout.GET("/events_for_month", r.handler.GetEventsForMonth)
}

func (r *Router) GetHTTPHandler() *gin.Engine {
	return r.rout
}
