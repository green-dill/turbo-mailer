package dashboard

import (
	"net/http"

	"turbo-mailer-server/internal/query"

	"github.com/labstack/echo/v4"
)

type DashboardStats struct {
	PoolCount      int64            `json:"pool_count"`
	TopPools       []TopPool        `json:"top_pools"`
	TaskCount      int64            `json:"task_count"`
	TaskStateCount map[string]int64 `json:"task_state_count"`
}

type TopPool struct {
	Name        string `json:"name"`
	SenderCount int    `json:"sender_count"`
}

// Stats godoc
//
//	@Summary		Get dashboard statistics
//	@Description	Retrieve statistics for the admin dashboard
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	DashboardStats
//	@Router			/api/v1/dashboard/stats [get]
func Stats(c echo.Context) error {
	ctx := c.Request().Context()

	var stats DashboardStats

	// Get pool count
	stats.PoolCount, _ = query.Pool.WithContext(ctx).Count()

	// Get top 10 pools by sender count
	topPools, _ := query.Pool.WithContext(ctx).Select(query.Pool.Name, query.Pool.SenderCount).
		Order(query.Pool.SenderCount.Desc()).
		Limit(10).
		Find()

	stats.TopPools = make([]TopPool, len(topPools))
	for i, pool := range topPools {
		stats.TopPools[i] = TopPool{
			Name:        pool.Name,
			SenderCount: pool.SenderCount,
		}
	}

	// Get task count
	stats.TaskCount, _ = query.Task.WithContext(ctx).Count()

	// Get task state count
	var taskStateCounts []struct {
		State string
		Count int64
	}
	query.Task.WithContext(ctx).Select(query.Task.State, query.Task.State.Count().As("count")).
		Group(query.Task.State).
		Scan(&taskStateCounts)

	stats.TaskStateCount = make(map[string]int64)
	for _, stateCount := range taskStateCounts {
		stats.TaskStateCount[stateCount.State] = stateCount.Count
	}

	return c.JSON(http.StatusOK, stats)
}
