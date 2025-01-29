package api

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/mburaksoran/insider-case/internal/infra/engines"
	"time"
)

type ServiceHealth struct {
	Name   string `json:"Name"`
	Status string `json:"Status"`
}

func checkRedis() string {
	redisClient := engines.GetRedisEngine()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := redisClient.Client.Ping(ctx).Result()
	if err != nil {
		return "unhealthy"
	}
	return "healthy"
}

// PostgreSQL Healthcheck
func checkPostgres() string {
	sqlClient := engines.GetSqlDbEngine()
	err := sqlClient.Client.Ping()
	if err != nil {
		return "unhealthy"
	}
	return "healthy"
}

// Vault Healthcheck
func checkVault() string {
	vaultClient := engines.GetVaultEngine()
	healthStatus, err := vaultClient.Client.Sys().Health()
	if err != nil || healthStatus.Sealed {
		return "unhealthy"
	}
	return "healthy"
}

func checkHealth() map[string]ServiceHealth {
	services := map[string]ServiceHealth{
		"redis":    {"redis", checkRedis()},
		"postgres": {"postgres", checkPostgres()},
		"vault":    {"vault", checkVault()},
	}
	return services
}
func HealthcheckHandler(c *fiber.Ctx) error {
	services := checkHealth()
	for _, service := range services {
		if service.Status == "unhealthy" {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":   "unhealthy",
				"services": services,
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":   "healthy",
		"services": services,
	})
}
