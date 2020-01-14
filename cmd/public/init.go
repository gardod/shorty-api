package public

import (
	"github.com/gardod/shorty-api/internal/driver/postgres"
	"github.com/gardod/shorty-api/internal/driver/redis"
)

func initDrivers() {
	postgres.InitDB()
	redis.InitClient()
}
