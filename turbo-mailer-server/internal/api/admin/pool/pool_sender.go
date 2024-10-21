package pool

import (
	"net/http"
	"strconv"
	"turbo-mailer-server/internal/models"
	"turbo-mailer-server/internal/query"
	"turbo-mailer-server/internal/schema"

	"github.com/labstack/echo/v4"
)

// SenderList
//
//	@Summary		List all pool senders for a specific pool
//	@Description	Get a paginated list of pool senders for a specific pool with optional sorting and filtering
//	@Tags			Pool Senders
//	@Accept			json
//	@Produce		json
//	@Param			id			path	int		true	"Pool ID"
//	@Param			page		query	int		false	"Page number"					default(1)
//	@Param			page_size	query	int		false	"Page size"						default(10)
//	@Param			sort		query	string	false	"Sort order: 'asc' or 'desc'"	default(desc)
//	@Param			search		query	string	false	"Search by from_email"
//	@Security		JWT
//	@Success		200	{object}	schema.Page[models.PoolSender]
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/pool-senders/{id} [get]
func SenderList(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse pool ID
	poolID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Pool ID"})
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize < 1 {
		pageSize = 10
	}
	sort := c.QueryParam("sort")
	search := c.QueryParam("search")

	// Build query
	q := query.PoolSender.WithContext(ctx).Where(query.PoolSender.PoolID.Eq(uint(poolID)))

	// Apply search filter if provided
	if search != "" {
		q = q.Where(query.PoolSender.FromEmail.Like("%" + search + "%"))
	}

	// Count total records
	total, err := q.Count()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Apply sorting
	if sort == "asc" {
		q = q.Order(query.PoolSender.ID)
	} else {
		q = q.Order(query.PoolSender.ID.Desc())
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	q = q.Offset(offset).Limit(pageSize)

	// Execute query
	poolSenders, err := q.Find()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Build next and prev URLs
	baseURL := c.Request().URL.Path + "?"
	nextPage := page + 1
	prevPage := page - 1
	next := ""
	prev := ""

	if int64(nextPage*pageSize) <= total {
		next = baseURL + "page=" + strconv.Itoa(nextPage) + "&page_size=" + strconv.Itoa(pageSize)
		if sort != "" {
			next += "&sort=" + sort
		}
		if search != "" {
			next += "&search=" + search
		}
	}

	if prevPage > 0 {
		prev = baseURL + "page=" + strconv.Itoa(prevPage) + "&page_size=" + strconv.Itoa(pageSize)
		if sort != "" {
			prev += "&sort=" + sort
		}
		if search != "" {
			prev += "&search=" + search
		}
	}

	// Create page response
	pageResponse := schema.NewPage(total, page, pageSize, poolSenders, next, prev)

	return c.JSON(http.StatusOK, pageResponse)
}

// SenderStore
//
//	@Summary		Create a new pool sender for a specific pool
//	@Description	Create a new pool sender with the provided details for a specific pool
//	@Tags			Pool Senders
//	@Accept			json
//	@Produce		json
//	@Param			id			path	int					true	"Pool ID"
//	@Param			poolSender	body	models.PoolSender	true	"Pool Sender object"
//	@Security		JWT
//	@Success		201	{object}	models.PoolSender
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/pool-senders/{id} [post]
func SenderStore(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse pool ID
	poolID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Pool ID"})
	}

	// Check if the pool exists
	pool, err := query.Pool.WithContext(ctx).Where(query.Pool.ID.Eq(uint(poolID))).First()
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Pool not found"})
	}

	// Pool exists, continue with the rest of the function

	poolSender := new(models.PoolSender)
	if err := c.Bind(poolSender); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Set the PoolID from the path parameter
	poolSender.PoolID = pool.ID

	// Validate required fields
	if poolSender.FromName == "" || poolSender.FromEmail == "" || poolSender.Domain == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "FromName, FromEmail, and Domain are required"})
	}

	err = query.PoolSender.WithContext(ctx).Create(poolSender)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, poolSender)
}

// SenderDelete
//
//	@Summary		Delete a pool sender
//	@Description	Delete a pool sender by ID
//	@Tags			Pool Senders
//	@Accept			json
//	@Produce		json
//	@Param			sid	path	int	true	"Pool Sender ID"
//	@Security		JWT
//	@Success		204	"No Content"
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/pool-senders/{sid} [delete]
func SenderDelete(c echo.Context) error {
	ctx := c.Request().Context()
	senderID, err := strconv.ParseUint(c.Param("sid"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Sender ID"})
	}

	_, err = query.PoolSender.WithContext(ctx).Where(query.PoolSender.ID.Eq(uint(senderID))).Unscoped().Delete()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
