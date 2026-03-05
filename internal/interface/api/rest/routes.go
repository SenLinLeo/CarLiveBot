package rest

import (
	"context"
	"net/http"
	"time"

	"carlivebot/internal/infrastructure/config"
	"carlivebot/internal/interface/api/rest/dto/request"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes 注册 HTTP 路由
func RegisterRoutes(e *echo.Echo, d *Deps) {
	cfg := d.Config
	repo := config.NewStoreConfigRepo(cfg.Server.ConfigPath)

	e.GET("/health", healthHandler())
	e.GET("/api/v1/live/status", liveStatusHandler())
	e.GET("/api/v1/stores", listStoresHandler(repo))
	e.GET("/api/v1/stores/:id", getStoreHandler(repo))
	e.GET("/api/v1/compliance/label", complianceLabelHandler(cfg))
	e.POST("/api/v1/script/generate", scriptGenerateHandler(d))
	e.GET("/api/v1/leads", leadsListHandler())
}

func healthHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	}
}

func liveStatusHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	}
}

func listStoresHandler(repo *config.StoreConfigRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		ids, err := repo.ListStoreIDs()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{"store_ids": ids})
	}
}

func getStoreHandler(repo *config.StoreConfigRepo) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "store id required"})
		}
		store, err := repo.GetByID(id)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, store)
	}
}

func complianceLabelHandler(cfg *config.App) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"label": cfg.Compliance.LabelText})
	}
}

func scriptGenerateHandler(d *Deps) echo.HandlerFunc {
	return func(c echo.Context) error {
		if d.Script == nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": "script service not configured"})
		}
		var req request.ScriptGenerateRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
		}
		if req.StoreID == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "store_id required"})
		}
		ctx, cancel := context.WithTimeout(c.Request().Context(), 60*time.Second)
		defer cancel()
		var fullText string
		err := d.Script.GenerateOnly(ctx, req.StoreID, req.UserInput, func(text string) error {
			fullText += text
			return nil
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]string{"text": fullText})
	}
}

func leadsListHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		storeID := c.QueryParam("store_id")
		_ = storeID
		return c.JSON(http.StatusOK, map[string]interface{}{
			"leads": []interface{}{},
			"total": 0,
		})
	}
}
