package handlers

import (
	"awesomeProject/internal/models"
	"awesomeProject/internal/service"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type CalendarHandler struct {
	calendarService *service.CalendarService
}

func NewCalendarHandler(calendarService *service.CalendarService) *CalendarHandler {
	return &CalendarHandler{calendarService: calendarService}
}

func (h *CalendarHandler) CreateEvent(c *gin.Context) {
	log := c.Value("logger").(*zap.Logger)
	log.Info("CreateEvent handler called")

	req := &models.EventRequest{}

	if err := json.NewDecoder(c.Request.Body).Decode(req); err != nil {
		log.Error("Failed to decode request body", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid request body"}) // Какой код возвращать?
		return
	}
	if req.UserID <= 0 || req.Event == "" || req.Date == "" {
		log.Error("Missing required parameters", zap.Int64("user_id", req.UserID), zap.String("date", req.Date), zap.String("event", req.Event))
		c.JSON(400, gin.H{"error": "Missing required parameters"})
		return
	}

	log.Info("Received CreateEvent request", zap.Int64("user_id", req.UserID), zap.String("date", req.Date), zap.String("event", req.Event))

	dateTime, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		log.Error("Invalid date format", zap.String("date", req.Date), zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}
	serviceEvent := &models.Event{
		UserID: req.UserID,
		Date:   dateTime,
		Event:  req.Event,
	}
	err = h.calendarService.CreateEvent(c.Request.Context(), serviceEvent)
	if err != nil {
		log.Error("Failed to create event", zap.Error(err))
		c.JSON(503, gin.H{"error": "Failed to create event"})
		return
	}
	log.Info("Event created successfully", zap.Int64("user_id", serviceEvent.UserID), zap.Time("date", dateTime), zap.String("event", serviceEvent.Event))
	c.JSON(200, gin.H{"result": "Event created successfully"})
}
func (h *CalendarHandler) UpdateEvent(c *gin.Context) {
	log := c.Value("logger").(*zap.Logger)
	log.Info("UpdateEvent handler called")
	req := &models.EventRequest{}

	if err := json.NewDecoder(c.Request.Body).Decode(req); err != nil {
		log.Error("Failed to decode request body", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid request body"}) // Какой код возвращать?
		return
	}
	if req.UserID <= 0 || req.Event == "" || req.Date == "" {
		log.Error("Missing required parameters", zap.Int64("user_id", req.UserID), zap.String("date", req.Date), zap.String("event", req.Event))
		c.JSON(400, gin.H{"error": "Missing required parameters"})
		return
	}

	log.Info("Received UpdateEvent request", zap.Int64("user_id", req.UserID), zap.String("date", req.Date), zap.String("event", req.Event))

	dateTime, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		log.Error("Invalid date format", zap.String("date", req.Date), zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}
	serviceEvent := &models.Event{
		UserID: req.UserID,
		Date:   dateTime,
		Event:  req.Event,
	}
	err = h.calendarService.UpdateEvent(c.Request.Context(), serviceEvent)
	if err != nil {
		log.Error("Failed to update event", zap.Error(err))
		c.JSON(503, gin.H{"error": "Failed to update event"})
		return
	}
	log.Info("Event updated successfully", zap.Int64("user_id", serviceEvent.UserID), zap.Time("date", dateTime), zap.String("event", serviceEvent.Event))
	c.JSON(200, gin.H{"result": "Event updated successfully"})
}
func (h *CalendarHandler) DeleteEvent(c *gin.Context) {
	log := c.Value("logger").(*zap.Logger)
	log.Info("DeleteEvent handler called")
	req := &models.EventRequest{}

	if err := json.NewDecoder(c.Request.Body).Decode(req); err != nil {
		log.Error("Failed to decode request body", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid request body"}) // Какой код возвращать?
		return
	}

	if req.UserID <= 0 || req.Date == "" {
		log.Error("Missing required parameters", zap.Int64("user_id", req.UserID), zap.String("date", req.Date))
		c.JSON(400, gin.H{"error": "Missing required parameters"})
		return
	}

	log.Info("Received DeleteEvent request", zap.Int64("user_id", req.UserID), zap.String("date", req.Date))

	dateTime, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		log.Error("Invalid date format", zap.String("date", req.Date), zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}
	serviceEvent := &models.Event{
		UserID: req.UserID,
		Date:   dateTime,
	}
	err = h.calendarService.DeleteEvent(c.Request.Context(), serviceEvent)
	if err != nil {
		log.Error("Failed to delete event", zap.Error(err))
		c.JSON(503, gin.H{"error": "Failed to delete event"})
		return
	}
	log.Info("Event deleted successfully", zap.Int64("user_id", serviceEvent.UserID), zap.Time("date", dateTime))
	c.JSON(200, gin.H{"result": "Event deleted successfully"})
}

func (h *CalendarHandler) GetEventsForDay(c *gin.Context) {
	log := c.Value("logger").(*zap.Logger)
	log.Info("GetEventsForDay handler called")

	userID, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil || userID <= 0 {
		log.Error("Invalid or missing user_id", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid or missing user_id parameter"})
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		log.Error("Missing date parameter")
		c.JSON(400, gin.H{"error": "Missing date parameter"})
		return
	}

	log.Info("Received GetEventsForDay request", zap.Int64("user_id", userID), zap.String("date", dateStr))

	dateTime, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Error("Invalid date format", zap.String("date", dateStr), zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}
	serviceEvent := &models.Event{
		UserID: userID,
		Date:   dateTime,
	}
	events, err := h.calendarService.GetEventsForDay(c.Request.Context(), serviceEvent)
	if err != nil {
		log.Error("Failed to get events for day", zap.Error(err))
		c.JSON(503, gin.H{"error": "Failed to get events for day"})
		return
	}
	log.Info("Events retrieved successfully", zap.Int64("user_id", serviceEvent.UserID), zap.Time("date", dateTime), zap.Int("event_count", len(events)))
	c.JSON(200, gin.H{"result": events})
}

func (h *CalendarHandler) GetEventsForWeek(c *gin.Context) {
	log := c.Value("logger").(*zap.Logger)
	log.Info("GetEventsForWeek handler called")
	userID, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil || userID <= 0 {
		log.Error("Invalid or missing user_id", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid or missing user_id parameter"})
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		log.Error("Missing date parameter")
		c.JSON(400, gin.H{"error": "Missing date parameter"})
		return
	}

	log.Info("Received GetEventsForWeek request", zap.Int64("user_id", userID), zap.String("date", dateStr))

	dateTime, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Error("Invalid date format", zap.String("date", dateStr), zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}
	serviceEvent := &models.Event{
		UserID: userID,
		Date:   dateTime,
	}
	events, err := h.calendarService.GetEventsForWeek(c.Request.Context(), serviceEvent)
	if err != nil {
		log.Error("Failed to get events for week", zap.Error(err))
		c.JSON(503, gin.H{"error": "Failed to get events for week"})
		return
	}
	log.Info("Events retrieved successfully", zap.Int64("user_id", serviceEvent.UserID), zap.Time("date", dateTime), zap.Int("event_count", len(events)))
	c.JSON(200, gin.H{"result": events})
}
func (h *CalendarHandler) GetEventsForMonth(c *gin.Context) {
	log := c.Value("logger").(*zap.Logger)
	log.Info("GetEventsForMonth handler called")
	userID, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil || userID <= 0 {
		log.Error("Invalid or missing user_id", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid or missing user_id parameter"})
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		log.Error("Missing date parameter")
		c.JSON(400, gin.H{"error": "Missing date parameter"})
		return
	}

	log.Info("Received GetEventsForMonth request", zap.Int64("user_id", userID), zap.String("date", dateStr))

	dateTime, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Error("Invalid date format", zap.String("date", dateStr), zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}
	serviceEvent := &models.Event{
		UserID: userID,
		Date:   dateTime,
	}
	events, err := h.calendarService.GetEventsForMonth(c.Request.Context(), serviceEvent)
	if err != nil {
		log.Error("Failed to get events for month", zap.Error(err))
		c.JSON(503, gin.H{"error": "Failed to get events for month"})
		return
	}
	log.Info("Events retrieved successfully", zap.Int64("user_id", serviceEvent.UserID), zap.Time("date", dateTime), zap.Int("event_count", len(events)))
	c.JSON(200, gin.H{"result": events})
}
