package public

import "github.com/gardod/shorty-api/internal/driver/postgres"

func initDrivers() {
	postgres.InitDB()
}
