package task

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"turbo-mailer-server/internal/dispatch"
	"turbo-mailer-server/internal/models"
	"turbo-mailer-server/internal/query"
	"turbo-mailer-server/internal/schema"
	"turbo-mailer-server/internal/utils/ptr"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

// List
//
//	@Summary		List all tasks with pagination
//	@Description	Get a paginated list of tasks with optional sorting and filtering
//	@Tags			Tasks
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int		false	"Page number"					default(1)
//	@Param			pageSize	query		int		false	"Page size"						default(10)
//	@Param			sort		query		string	false	"Sort order: 'asc' or 'desc'"	default(desc)
//	@Param			search		query		string	false	"Search by task subject"
//	@Success		200			{object}	schema.Page[models.Task]
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Security		JWT
//	@Router			/api/v1/tasks [get]
func List(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse query parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	if pageSize < 1 {
		pageSize = 10
	}
	sort := c.QueryParam("sort")
	search := c.QueryParam("search")

	// Build query
	q := query.Task.WithContext(ctx).Preload(query.Task.Pools)

	// Apply search filter if provided
	if search != "" {
		q = q.Where(query.Task.Subject.Like("%" + search + "%"))
	}

	// Count total records
	total, err := q.Count()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Apply sorting
	if sort == "asc" {
		q = q.Order(query.Task.ID)
	} else {
		q = q.Order(query.Task.ID.Desc())
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	q = q.Offset(offset).Limit(pageSize)

	// Execute query
	tasks, err := q.Find()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Limit the number of receivers to 100
	// reduce the response payload size
	for _, task := range tasks {
		if len(task.Receivers.Val()) > 100 {
			task.Receivers = models.NewJSON(task.Receivers.Val()[:100])
		}
	}

	// Build next and prev URLs
	baseURL := c.Request().URL.Path
	nextPage := page + 1
	prevPage := page - 1
	next := ""
	prev := ""

	if int64(nextPage*pageSize) <= total {
		nextValues := url.Values{}
		nextValues.Set("page", strconv.Itoa(nextPage))
		nextValues.Set("pageSize", strconv.Itoa(pageSize))
		if sort != "" {
			nextValues.Set("sort", sort)
		}
		if search != "" {
			nextValues.Set("search", search)
		}
		next = baseURL + "?" + nextValues.Encode()
	}

	if prevPage > 0 {
		prevValues := url.Values{}
		prevValues.Set("page", strconv.Itoa(prevPage))
		prevValues.Set("pageSize", strconv.Itoa(pageSize))
		if sort != "" {
			prevValues.Set("sort", sort)
		}
		if search != "" {
			prevValues.Set("search", search)
		}
		prev = baseURL + "?" + prevValues.Encode()
	}

	// Create page response
	pageResponse := schema.NewPage(total, page, pageSize, tasks, next, prev)

	return c.JSON(http.StatusOK, pageResponse)
}

// Get
//
//	@Summary		Get a specific task
//	@Description	Get details of a specific task by ID
//	@Tags			Tasks
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Task ID"
//	@Success		200	{object}	models.Task
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Security		JWT
//	@Router			/api/v1/tasks/{id} [get]
func Get(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	task, err := query.Task.WithContext(ctx).Preload(query.Task.Pools.Pool.Senders).Where(query.Task.ID.Eq(uint(id))).First()
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Task not found"})
	}

	// Limit the number of receivers to 100
	// reduce the response payload size
	if len(task.Receivers.Val()) > 100 {
		task.Receivers = models.NewJSON(task.Receivers.Val()[:100])
	}

	return c.JSON(http.StatusOK, task)
}

// Store
//
//	@Summary		Create a new task
//	@Description	Create a new task with the provided details and uploaded files
//	@Tags			Tasks
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			subject					formData	string	true	"Task subject"
//	@Param			context_type			formData	string	true	"Context type (html or text)"
//	@Param			content					formData	file	true	"Content file (html or text)"
//	@Param			receivers				formData	file	true	"Receivers CSV file"
//	@Param			max_dispatch_per_hour	formData	int		false	"Max dispatch per hour"
//	@Param			schedule_at				formData	string	false	"Schedule at, format: YYYY-MM-DD HH:MM:SS"
//	@Param			pools					formData	[]int	false	"Pools"
//	@Param			pools_weights			formData	[]int	false	"Pools weights"
//	@Param			metadata				formData	string	false	"Metadata, JSON string of map[string]string"
//	@Success		201						{object}	models.Task
//	@Failure		400						{object}	map[string]string
//	@Failure		500						{object}	map[string]string
//	@Security		JWT
//	@Router			/api/v1/tasks [post]
func Store(c echo.Context) error {
	ctx := c.Request().Context()
	task := new(models.Task)

	// Parse form data
	task.Subject = c.FormValue("subject")
	task.ContentType = c.FormValue("content_type")
	task.State = models.TaskStatePending
	maxDispatchPerHour, err := strconv.Atoi(c.FormValue("max_dispatch_per_hour"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid max_dispatch_per_hour", "detail": err.Error()})
	}
	task.MaxDispatchPreHour = maxDispatchPerHour

	// Handle content file
	contentFile, err := c.FormFile("content")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Content file is required", "detail": err.Error()})
	}
	content, err := readFile(contentFile)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read content file", "detail": err.Error()})
	}
	task.Content = string(content)

	// Handle receivers file
	receiversFile, err := c.FormFile("receivers")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Receivers file is required", "detail": err.Error()})
	}
	receivers, err := readCSV(receiversFile)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read receivers file", "detail": err.Error()})
	}
	task.Receivers = models.NewJSON(receivers)

	if scheduleAt := c.FormValue("schedule_at"); scheduleAt != "" {
		scheduleAtTime, err := time.Parse(time.DateTime, scheduleAt)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid schedule_at", "detail": err.Error()})
		}
		task.ScheduleAt = &scheduleAtTime
	} else {
		task.ScheduleAt = ptr.Ptr(time.Now())
	}
	formPools := c.FormValue("pools")
	formPoolsWeights := c.FormValue("pools_weights")

	if metadata := c.FormValue("metadata"); metadata != "" {
		metadataMap := make(map[string]string)
		raw := []byte(metadata)
		if err := json.Unmarshal(raw, &metadataMap); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid metadata", "detail": err.Error()})
		} else {
			task.Metadata = ptr.Ptr(json.RawMessage(raw))
		}
	}

	var (
		poolIds []uint
		weights []int
		pools   []*models.TaskPool
	)

	if formPools != "" {
		poolIds = lo.Map(strings.Split(formPools, ","), func(poolId string, _ int) uint {
			poolIdUint, err := strconv.ParseUint(poolId, 10, 32)
			if err != nil {
				return 0
			}
			return uint(poolIdUint)
		})
	}

	if formPoolsWeights != "" {
		weights = lo.Map(strings.Split(formPoolsWeights, ","), func(weight string, _ int) int {
			weightInt, err := strconv.Atoi(weight)
			if err != nil {
				return 0
			}
			return weightInt
		})
	}

	if len(poolIds) != len(weights) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Pools and weights must have the same length"})
	}

	// pool ids must exists
	ps, err := query.Pool.WithContext(ctx).
		Where(query.Pool.ID.In(poolIds...)).
		Where(query.Pool.SenderCount.Gt(0)).
		Find()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get pools", "detail": err.Error()})
	}
	if len(ps) != len(poolIds) {
		missingPoolIds := lo.Filter(poolIds, func(poolId uint, _ int) bool {
			return !lo.ContainsBy(ps, func(pool *models.Pool) bool {
				return pool.ID == poolId
			})
		})
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("some of the pools are not found: %v", missingPoolIds)})
	}

	// create task pools
	for i, p := range ps {
		taskPool := &models.TaskPool{
			PoolID: p.ID,
			Pool:   p,
			Weight: weights[i],
		}
		pools = append(pools, taskPool)
	}

	query.DB.Transaction(func(tx *gorm.DB) error {
		// err = query.TaskPool.WithContext(ctx).Save(pools...)
		// if err != nil {
		// 	return err
		// }
		task.Pools = pools
		err = query.Task.WithContext(ctx).Create(task)
		if err != nil {
			return err
		}
		return nil
	})

	// Limit the number of receivers to 100
	// reduce the response payload size
	if len(receivers) > 100 {
		task.Receivers = models.NewJSON(receivers[:100])
	}

	return c.JSON(http.StatusCreated, task)
}

// Update
//
//	@Summary		Update a task
//	@Description	Update an existing task's details and files
//	@Tags			Tasks
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id						path		int		true	"Task ID"
//	@Param			subject					formData	string	false	"Task subject"
//	@Param			context_type			formData	string	false	"Context type (html or text)"
//	@Param			content					formData	file	false	"Content file (html or text)"
//	@Param			receivers				formData	file	false	"Receivers CSV file"
//	@Param			state					formData	string	false	"Task state"
//	@Param			max_dispatch_per_hour	formData	int		false	"Max dispatch per hour"
//	@Param			schedule_at				formData	string	false	"Schedule at, format: YYYY-MM-DD HH:MM:SS"
//	@Param			metadata				formData	string	false	"Metadata, JSON string of map[string]string"
//	@Param			pools					formData	[]int	false	"Pools"
//	@Param			pools_weights			formData	[]int	false	"Pools weights"
//	@Success		200						{object}	models.Task
//	@Failure		400						{object}	map[string]string
//	@Failure		404						{object}	map[string]string
//	@Failure		500						{object}	map[string]string
//	@Security		JWT
//	@Router			/api/v1/tasks/{id} [put]
func Update(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	// Fetch the existing task
	existingTask, err := query.Task.WithContext(ctx).Where(query.Task.ID.Eq(uint(id))).First()
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Task not found"})
	}

	// Update fields if provided
	if subject := c.FormValue("subject"); subject != "" {
		existingTask.Subject = subject
	}
	if contextType := c.FormValue("context_type"); contextType != "" {
		existingTask.ContentType = contextType
	}
	if state := c.FormValue("state"); state != "" {
		existingTask.State = state
	}
	if maxDispatchPerHourStr := c.FormValue("max_dispatch_per_hour"); maxDispatchPerHourStr != "" {
		maxDispatchPerHour, err := strconv.Atoi(maxDispatchPerHourStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid maxDispatchPerHour"})
		}
		existingTask.MaxDispatchPreHour = maxDispatchPerHour
	}

	if metadata := c.FormValue("metadata"); metadata != "" {
		metadataMap := make(map[string]string)
		raw := []byte(metadata)
		if err := json.Unmarshal(raw, &metadataMap); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid metadata"})
		} else {
			existingTask.Metadata = ptr.Ptr(json.RawMessage(raw))
		}
	}

	// Handle content file if provided
	if contentFile, err := c.FormFile("content"); err == nil {
		content, err := readFile(contentFile)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read content file"})
		}
		existingTask.Content = string(content)
	}

	// Handle receivers file if provided
	if receiversFile, err := c.FormFile("receivers"); err == nil {
		receivers, err := readCSV(receiversFile)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read receivers file"})
		}
		existingTask.Receivers = models.NewJSON(receivers)
	}

	if scheduleAt := c.FormValue("schedule_at"); scheduleAt != "" {
		scheduleAtTime, err := time.Parse(time.DateTime, scheduleAt)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid schedule_at"})
		}
		existingTask.ScheduleAt = &scheduleAtTime
	}

	formPools := c.FormValue("pools")
	formPoolsWeights := c.FormValue("pools_weights")

	var (
		poolIds []uint
		weights []int
		pools   []*models.TaskPool
	)

	if formPools != "" {
		poolIds = lo.Map(strings.Split(formPools, ","), func(poolId string, _ int) uint {
			poolIdUint, err := strconv.ParseUint(poolId, 10, 32)
			if err != nil {
				return 0
			}
			return uint(poolIdUint)
		})
	}

	if formPoolsWeights != "" {
		weights = lo.Map(strings.Split(formPoolsWeights, ","), func(weight string, _ int) int {
			weightInt, err := strconv.Atoi(weight)
			if err != nil {
				return 0
			}
			return weightInt
		})
	}

	if len(poolIds) != len(weights) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Pools and weights must have the same length"})
	}

	ps, err := query.Pool.WithContext(ctx).
		Select(query.Pool.ID, query.Pool.Name).
		Where(query.Pool.ID.In(poolIds...)).
		Where(query.Pool.SenderCount.Gt(0)).
		Find()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get pools"})
	}
	if len(ps) != len(poolIds) {
		missingPoolIds := lo.Filter(poolIds, func(poolId uint, _ int) bool {
			return !lo.ContainsBy(ps, func(pool *models.Pool) bool {
				return pool.ID == poolId
			})
		})
		return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("some of the pools are not found: %v", missingPoolIds)})
	}

	for i, poolId := range poolIds {
		pools = append(pools, &models.TaskPool{
			PoolID: poolId,
			Pool:   ps[i],
			Weight: weights[i],
		})
	}

	// clean up existing pools
	if len(existingTask.Pools) > 0 {
		query.TaskPool.WithContext(ctx).Where(query.TaskPool.TaskID.Eq(existingTask.ID)).Unscoped().Delete()
	}

	existingTask.Pools = pools

	// Perform the update
	_, err = query.Task.WithContext(ctx).Where(query.Task.ID.Eq(uint(id))).Updates(existingTask)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, existingTask)
}

// Delete
//
//	@Summary		Delete a task
//	@Description	Delete a task by ID
//	@Tags			Tasks
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Task ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Security		JWT
//	@Router			/api/v1/tasks/{id} [delete]
func Delete(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	task, err := query.Task.WithContext(ctx).Preload(query.Task.Pools).Where(query.Task.ID.Eq(uint(id))).First()
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Task not found"})
	}

	if task.State == models.TaskStateDispatched {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Task is sending, cannot delete"})
	}

	_, err = query.Task.WithContext(ctx).Where(query.Task.ID.Eq(uint(id))).Delete()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// Test
//
//	@Summary		Test a task
//	@Description	Test a task by sending a test email to the specified address
//	@Tags			Tasks
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int		true	"Task ID"
//	@Param			email	query		string	true	"Email address to send the test to"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Security		JWT
//	@Router			/api/v1/tasks/{id}/test [post]
func Test(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid task ID"})
	}

	email := c.QueryParam("email")
	if email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Test email is required"})
	}

	// Fetch the task
	task, err := query.Task.WithContext(ctx).Where(query.Task.ID.Eq(uint(id))).First()
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Task not found"})
	}

	if task.State != models.TaskStatePending {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Task is not pending"})
	}

	err = dispatch.TestSend(ctx, task.ID, email)
	if err != nil {
		log.Error().Err(err).Msg("failed to send test email")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to send test email: %s", err)})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Test initiated",
		"task_id": strconv.FormatUint(uint64(task.ID), 10),
		"email":   email,
	})
}

// StartImmediately
//
//	@Summary		Start a task immediately
//	@Description	Immediately start a task by setting its schedule time to now
//	@Tags			Tasks
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Task ID"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Security		JWT
//	@Router			/api/v1/tasks/{id}/start-immediately [post]
func StartImmediately(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid task ID"})
	}
	// Fetch the task
	task, err := query.Task.WithContext(ctx).Where(query.Task.ID.Eq(uint(id))).First()
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Task not found"})
	}

	if task.State != models.TaskStatePending {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Task is not pending"})
	}

	// Update the schedule_at to current time
	now := time.Now()
	_, err = query.Task.WithContext(ctx).Where(query.Task.ID.Eq(uint(id))).Update(query.Task.ScheduleAt, now)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update task schedule"})
	}

	err = dispatch.DispatchImmediately(ctx, task)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to dispatch task"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message":    "Task scheduled to start immediately",
		"taskId":     strconv.FormatUint(uint64(task.ID), 10),
		"scheduleAt": now.Format(time.RFC3339),
	})

}

// DownloadContentTemplate
//
//	@Summary		Download content template
//	@Description	Download a template file for the task content (HTML or plain text)
//	@Tags			Tasks
//	@Accept			json
//	@Produce		octet-stream
//	@Param			type	query		string	true	"Template type (html or text)"
//	@Success		200		{file}		binary	"Template file"
//	@Failure		400		{object}	map[string]string
//	@Security		JWT
//	@Router			/api/v1/tasks/content-template [get]
func DownloadContentTemplate(c echo.Context) error {
	templateType := c.QueryParam("type")
	if templateType != "html" && templateType != "text" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid template type"})
	}

	var filename string
	var content string

	if templateType == "html" {
		filename = "content_template.html"
		content = `<!DOCTYPE html>
<html>
<head>
    <title>Email Template</title>
</head>
<body>
    <h1>Hello, {{name}}!</h1>
    <p>This is a sample email template.</p>
</body>
</html>`
	} else {
		filename = "content_template.txt"
		content = `Hello, {{name}}!

This is a sample email template.`
	}

	c.Response().Header().Set("Content-Disposition", "attachment; filename="+filename)
	return c.Blob(http.StatusOK, echo.MIMEOctetStream, []byte(content))
}

// DownloadReceiversTemplate
//
//	@Summary		Download receivers template
//	@Description	Download a template CSV file for the task receivers with example email addresses
//	@Tags			Tasks
//	@Accept			json
//	@Produce		text/csv
//	@Success		200	{file}	binary	"CSV template file"
//	@Security		JWT
//	@Router			/api/v1/tasks/receivers-template [get]
func DownloadReceiversTemplate(c echo.Context) error {
	content := `email
john@example.com
jane@example.com
user@example.com`

	c.Response().Header().Set("Content-Disposition", "attachment; filename=receivers_template.csv")
	return c.Blob(http.StatusOK, echo.MIMEOctetStream, []byte(content))
}

// Helper functions

func readFile(file *multipart.FileHeader) ([]byte, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	return io.ReadAll(src)
}

func readCSV(file *multipart.FileHeader) ([]string, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	reader := csv.NewReader(src)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var emails []string
	for i, record := range records {
		if i == 0 || len(record) == 0 { // Skip header and empty rows
			continue
		}
		emails = append(emails, record[0])
	}

	return emails, nil
}
