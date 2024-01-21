package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"infotecs_trainee_task/internal/service"
	"log/slog"
	"os"
)

func NewRouter(handler *echo.Echo, services *service.Services) {
	handler.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}", "method":"${method}","uri":"${uri}", "status":${status},"error":"${error}"}` + "\n",
		Output: setLogsFile(),
	}))

	handler.Use(middleware.Recover())

	handler.GET("/health", func(c echo.Context) error {
		return c.NoContent(200)
	})

	auth := handler.Group("/auth")
	{
		newAuthRoutes(auth, services.Auth)
	}

	authMiddleware := &AuthMiddleware{services.Auth}
	v1 := handler.Group("/api/v1", authMiddleware.UserIdentity)
	{
		newWalletRoutes(v1.Group("/wallet"), services.Wallet)
	}
}

func setLogsFile() *os.File {
	file, err := os.OpenFile("/logs/request.log", os.O_CREATE|os.O_APPEND|os.O_APPEND, 0666)

	if err != nil {
		slog.Error("can not open file", err)
	}

	return file
}
