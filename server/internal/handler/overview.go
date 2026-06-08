package handler

import (
	"todo-server/internal/model"
	"todo-server/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OverviewHandler struct {
	overviewService *service.OverviewService
	taskService     *service.TaskService
}

func NewOverviewHandler(overviewService *service.OverviewService, taskService *service.TaskService) *OverviewHandler {
	return &OverviewHandler{
		overviewService: overviewService,
		taskService:     taskService,
	}
}

// GetOverview 获取仪表盘概览（聚合查询）
func (h *OverviewHandler) GetOverview(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	overview, err := h.overviewService.GetOverview(userID)
	if err != nil {
		model.Error(c, model.CodeInternalError, err.Error())
		return
	}

	model.Success(c, overview)
}

// GetTodayTasks 获取今日任务
func (h *OverviewHandler) GetTodayTasks(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	tasks, err := h.taskService.GetTodayTasks(userID)
	if err != nil {
		model.Error(c, model.CodeInternalError, err.Error())
		return
	}

	model.Success(c, tasks)
}

// GetUpcomingTasks 获取即将到期任务
func (h *OverviewHandler) GetUpcomingTasks(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	tasks, err := h.taskService.GetUpcomingTasks(userID)
	if err != nil {
		model.Error(c, model.CodeInternalError, err.Error())
		return
	}

	model.Success(c, tasks)
}
