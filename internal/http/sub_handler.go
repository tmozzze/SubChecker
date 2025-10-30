package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tmozzze/SubChecker/internal/model"
	"github.com/tmozzze/SubChecker/internal/service"
)

type SubHandler struct {
	svc service.SubService
	log *logrus.Logger
}

func NewSubHandler(s service.SubService, l *logrus.Logger) *SubHandler {
	return &SubHandler{svc: s, log: l}
}

type createSubReq struct {
	ServiceName string `json:"service_name" binding:"required"`
	Price       int    `json:"price" binding:"required,min=0"`
	UserId      string `json:"user_id" binding:"required,uuid"`
	StartDate   string `json:"start_date" binding:"required"` // MM-YYYY
	EndDate     string `json:"end_date,omitempty"`            // Optional
}

func parseMonth(s string) (time.Time, error) {
	t, err := time.Parse("01-2006", s)
	if err == nil {
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC), nil
	}
	return time.Time{}, err
}

// CreateSub godoc
// @Summary Create subscription
// @Description Create subscription record
// @Tags subs
// @Accept json
// @Produce json
// @Param body body createSubReq true "Subscription"
// @Success 201 {object} model.Sub
// @Failure 400 {object} gin.H
// @Router /subs [post]
func (h *SubHandler) CreateSub(c *gin.Context) {
	var req createSubReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithError(err).Warn("invalid create request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sd, err := parseMonth(req.StartDate)
	if err != nil {
		h.log.WithError(err).Warn("invalid parse month for start_date")
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad start_date format, expected MM-YYYY"})
		return
	}

	var ed *time.Time
	if req.EndDate != "" {
		t, err := parseMonth(req.EndDate)
		if err != nil {
			h.log.WithError(err).Warn("invalid create request for end_date")
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad end_date format, expected MM-YYYY"})
			return
		}
		ed = &t
	}

	sub := &model.Sub{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserId:      req.UserId,
		StartDate:   sd,
		EndDate:     ed,
	}
	if err := h.svc.Create(c.Request.Context(), sub); err != nil {
		h.log.WithError(err).Error("failed create sub")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusCreated, sub)
}

type sumReq struct {
	UserId      string `form:"user_id"`
	ServiceName string `form:"service_name"`
	StartMonth  string `form:"start_month" binding:"required"` // MM-YYYY
	EndMonth    string `form:"end_month" binding:"required"`   // MM-YYYY
}

// @Summary Sum cost
// @Description Sum total cost for period (inclusive months). Filters: user_id, service_name
// @Tags subs
// @Accept json
// @Produce json
// @Param start_month query string true "MM-YYYY"
// @Param end_month query string true "MM-YYYY"
// @Param user_id query string false "UUID"
// @Param service_name query string false "service name"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router /subs/sum [get]
func (h *SubHandler) SumCost(c *gin.Context) {
	var q sumReq
	if err := c.ShouldBindQuery(&q); err != nil {
		h.log.WithError(err).Warn("invalid create request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pStart, err := parseMonth(q.StartMonth)
	if err != nil {
		h.log.WithError(err).Warn("invalid parse month for start_month")
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad start_month"})
		return
	}
	pEnd, err := parseMonth(q.EndMonth)
	if err != nil {
		h.log.WithError(err).Warn("invalid parse month for end_month")
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad end_month"})
		return
	}

	total, err := h.svc.SumCost(c.Request.Context(), q.UserId, q.ServiceName, pStart, pEnd)
	if err != nil {
		h.log.WithError(err).Error("sum cost failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total_rub": total})
}

// GetSubByID godoc
// @Summary Get subscription by ID
// @Description Get subscription by its ID
// @Tags subs
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} model.Sub
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /subs/{id} [get]
func (h *SubHandler) GetSubById(c *gin.Context) {
	idStr := c.Param("sub_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.WithError(err).Warn("invalid sub_id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sub_id"})
		return
	}

	sub, err := h.svc.GetById(c.Request.Context(), id)
	if err != nil {
		h.log.WithError(err).Warn("get by id failed")
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	c.JSON(http.StatusOK, sub)

}

// ListSubs godoc
// @Summary List subscriptions
// @Description Get paginated list of subscriptions
// @Tags subs
// @Produce json
// @Param limit query int false "limit (default 50)"
// @Param offset query int false "offset (default 0)"
// @Success 200 {array} model.Sub
// @Router /subs [get]
func (h *SubHandler) ListSubs(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	subs, err := h.svc.List(c.Request.Context(), limit, offset)
	if err != nil {
		h.log.WithError(err).Error("list subs failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}

	c.JSON(http.StatusOK, subs)
}

// UpdateSub godoc
// @Summary Update subscription
// @Description Update existing subscription by ID
// @Tags subs
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Param body body createSubReq true "Updated subscription"
// @Success 200 {object} model.Sub
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /subs/{id} [put]
func (h *SubHandler) UpdateSub(c *gin.Context) {
	idStr := c.Param("sub_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sub_id"})
		return
	}

	var req createSubReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.WithError(err).Warn("invalid update request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sd, err := parseMonth(req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad start_date format, expected MM-YYYY"})
		return
	}

	var ed *time.Time
	if req.EndDate != "" {
		t, err := parseMonth(req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad end_date format, expected MM-YYYY"})
			return
		}
		ed = &t
	}

	sub := &model.Sub{
		SubId:       id,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserId:      req.UserId,
		StartDate:   sd,
		EndDate:     ed,
	}

	if err := h.svc.Update(c.Request.Context(), sub); err != nil {
		h.log.WithError(err).Error("update failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}

	c.JSON(http.StatusOK, sub)
}

// DeleteSub godoc
// @Summary Delete subscription
// @Description Delete subscription by ID
// @Tags subs
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 204 {object} nil
// @Failure 400 {object} gin.H
// @Router /subs/{id} [delete]
func (h *SubHandler) DeleteSub(c *gin.Context) {
	idStr := c.Param("sub_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sub_id"})
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		h.log.WithError(err).Error("delete failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}

	c.Status(http.StatusNoContent)
}
