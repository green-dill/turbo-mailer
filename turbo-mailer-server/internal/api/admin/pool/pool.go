package pool

import (
	"net/http"
	"strconv"
	"strings"
	"turbo-mailer-server/internal/models"
	"turbo-mailer-server/internal/query"
	"turbo-mailer-server/internal/schema"

	"github.com/labstack/echo/v4"
)

// List
//
//	@Summary		List all pools with pagination
//	@Description	Get a paginated list of pools with optional sorting and filtering
//	@Tags			Pools
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int		false	"Page number"					default(1)
//	@Param			page_size	query		int		false	"Page size"						default(10)
//	@Param			sort		query		string	false	"Sort order: 'asc' or 'desc'"	default(desc)
//	@Param			search		query		string	false	"Search by pool name"
//	@Param			has_sender	query		boolean	false	"Has sender"
//	@Success		200			{object}	schema.Page[models.Pool]
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Security		JWT
//	@Router			/api/v1/pool [get]
func List(c echo.Context) error {
	ctx := c.Request().Context()

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
	hasSender := c.QueryParam("has_sender")

	// Build query
	q := query.Pool.WithContext(ctx).Preload(query.Pool.Senders)

	// Apply search filter if provided
	if search != "" {
		q = q.Where(query.Pool.Name.Like("%" + search + "%"))
	}

	// Apply has sender filter if provided
	if strings.EqualFold(hasSender, "true") {
		q = q.Where(query.Pool.SenderCount.Gt(0))
	}

	// Count total records
	total, err := q.Count()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Apply sorting
	if sort == "asc" {
		q = q.Order(query.Pool.ID)
	} else {
		q = q.Order(query.Pool.ID.Desc())
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	q = q.Offset(offset).Limit(pageSize)

	// Execute query
	pools, err := q.Find()
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
	pageResponse := schema.NewPage(total, page, pageSize, pools, next, prev)

	return c.JSON(http.StatusOK, pageResponse)
}

// Get
//
//	@Summary		Get a specific pool
//	@Description	Get details of a specific pool by ID
//	@Tags			Pools
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Pool ID"
//	@Success		200	{object}	models.Pool
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Security		JWT
//	@Router			/api/v1/pool/{id} [get]
func Get(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	pool, err := query.Pool.WithContext(ctx).Preload(query.Pool.Senders).Where(query.Pool.ID.Eq(uint(id))).First()
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Pool not found"})
	}

	return c.JSON(http.StatusOK, pool)
}

// Store
//
//	@Summary		Create a new pool
//	@Description	Create a new pool with the provided details
//	@Tags			Pools
//	@Accept			json
//	@Produce		json
//	@Param			pool	body		models.Pool	true	"Pool object"
//	@Success		201		{object}	models.Pool
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Security		JWT
//	@Router			/api/v1/pool [post]
func Store(c echo.Context) error {
	ctx := c.Request().Context()
	pool := new(models.Pool)
	if err := c.Bind(pool); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err := query.Pool.WithContext(ctx).Create(pool)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, pool)
}

type UpdatePool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Update
//
//	@Summary		Update a pool
//	@Description	Update an existing pool's details
//	@Tags			Pools
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int			true	"Pool ID"
//	@Param			pool	body		UpdatePool	true	"Updated pool object"
//	@Success		200		{object}	models.Pool
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Security		JWT
//	@Router			/api/v1/pool/{id} [put]
func Update(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	updatePool := new(UpdatePool)
	if err := c.Bind(updatePool); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Fetch the existing pool
	existingPool, err := query.Pool.WithContext(ctx).Where(query.Pool.ID.Eq(uint(id))).First()
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Pool not found"})
	}

	// Update fields
	if updatePool.Name != "" {
		existingPool.Name = updatePool.Name
	}
	existingPool.Description = updatePool.Description

	// Perform the update
	_, err = query.Pool.WithContext(ctx).Where(query.Pool.ID.Eq(uint(id))).Updates(existingPool)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, existingPool)
}

// Delete
//
//	@Summary		Delete a pool
//	@Description	Delete a pool by ID
//	@Tags			Pools
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Pool ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Security		JWT
//	@Router			/api/v1/pool/{id} [delete]
func Delete(c echo.Context) error {
	ctx := c.Request().Context()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	_, err = query.Pool.WithContext(ctx).Where(query.Pool.ID.Eq(uint(id))).Delete()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
