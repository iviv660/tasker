package handler

import (
	"app/internal/entity"
	"app/internal/usecase"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Handler struct {
	TaskUseCase *usecase.TaskUseCase
	UserUseCase *usecase.UserUseCase
}

func NewHandler(taskUC *usecase.TaskUseCase, userUC *usecase.UserUseCase) (*gin.Engine, *Handler) {
	h := &Handler{TaskUseCase: taskUC, UserUseCase: userUC}
	r := gin.New()
	r.Use(gin.Recovery())

	// Публичные
	r.POST("/auth/register", h.registerUser)
	r.POST("/auth/login", h.login)

	// Защищённые
	auth := r.Group("/")
	auth.Use(AuthMiddleware())
	{
		auth.POST("/tasks", h.createTask)                  // создать задачу
		auth.GET("/tasks", h.getTasks)                     // список моих задач
		auth.GET("/tasks/:id", h.getTaskByID)              // получить одну задачу
		auth.PUT("/tasks/:id", h.updateTask)               // обновить задачу
		auth.PATCH("/tasks/:id/complete", h.completedTask) // отметить выполненной
		auth.DELETE("/tasks/:id", h.deleteTask)            // удалить задачу
	}

	return r, h
}

// ===== helpers =====

func getUserID(c *gin.Context) (int64, bool) {
	v, ok := c.Get("user_id")
	if !ok {
		return 0, false
	}
	id, ok := v.(int64)
	return id, ok
}

func parseIDParam(c *gin.Context, name string) (int64, bool) {
	s := c.Param(name)
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + name})
		return 0, false
	}
	return id, true
}

// ===== auth =====

// @Summary      Регистрация
// @Description  Создаёт пользователя, хэширует пароль
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "payload"
// @Success      200 {object} map[string]int64 "user_id"
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /auth/register [post]
func (h *Handler) registerUser(c *gin.Context) {
	var r RegisterRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	id, err := h.UserUseCase.Register(c.Request.Context(), r.Email, r.Password, r.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user_id": id})
}

// @Summary      Логин
// @Description  Возвращает JWT при валидных email/пароле
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "payload"
// @Success      200 {object} map[string]string "token"
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /auth/login [post]
func (h *Handler) login(c *gin.Context) {
	var r LoginRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	token, err := h.UserUseCase.Login(c.Request.Context(), r.Email, r.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// ===== tasks =====

// @Summary      Создать задачу
// @Security     BearerAuth
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        request body CreateTaskRequest true "payload"
// @Success      200 {object} entity.Task
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /tasks [post]
func (h *Handler) createTask(c *gin.Context) {
	var r CreateTaskRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user in context"})
		return
	}
	task, err := h.TaskUseCase.CreateTask(c.Request.Context(), userID, r.Title, r.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// @Summary      Мои задачи
// @Security     BearerAuth
// @Tags         tasks
// @Produce      json
// @Success      200 {object} TasksResponse
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /tasks [get]
func (h *Handler) getTasks(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user in context"})
		return
	}
	tasks, err := h.TaskUseCase.ListTasks(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list tasks"})
		return
	}
	c.JSON(http.StatusOK, TasksResponse{Tasks: tasks})
}

// @Summary      Одна задача
// @Security     BearerAuth
// @Tags         tasks
// @Produce      json
// @Param        id   path int true "Task ID"
// @Success      200 {object} entity.Task
// @Failure      401 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /tasks/{id} [get]
func (h *Handler) getTaskByID(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user in context"})
		return
	}
	taskID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	task, err := h.TaskUseCase.GetTaskByID(c.Request.Context(), taskID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get task"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// @Summary      Обновить задачу
// @Security     BearerAuth
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id   path int true "Task ID"
// @Param        request body UpdateTaskRequest true "payload"
// @Success      200 {object} entity.Task
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /tasks/{id} [put]
func (h *Handler) updateTask(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user in context"})
		return
	}
	taskID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}
	var r UpdateTaskRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	task, err := h.TaskUseCase.UpdateTask(c.Request.Context(), &entity.Task{
		ID:          taskID,
		OwnerID:     userID,
		Title:       r.Title,
		Description: r.Description,
		Status:      r.Status,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// @Summary      Отметить выполненной
// @Security     BearerAuth
// @Tags         tasks
// @Produce      json
// @Param        id   path int true "Task ID"
// @Success      200 {object} entity.Task
// @Failure      401 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /tasks/{id}/complete [patch]
func (h *Handler) completedTask(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user in context"})
		return
	}
	taskID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	task, err := h.TaskUseCase.UpdateTask(c.Request.Context(), &entity.Task{
		ID:      taskID,
		OwnerID: userID,
		Status:  true,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to complete task"})
		return
	}
	c.JSON(http.StatusOK, task)
}

// @Summary      Удалить задачу
// @Security     BearerAuth
// @Tags         tasks
// @Param        id   path int true "Task ID"
// @Success      204  "no content"
// @Failure      401 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /tasks/{id} [delete]
func (h *Handler) deleteTask(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user in context"})
		return
	}
	taskID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	if err := h.TaskUseCase.DeleteTask(c.Request.Context(), taskID, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task"})
		return
	}
	c.Status(http.StatusNoContent)
}
