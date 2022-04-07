package leopard

import "os"

// GetEnv get a enviroment variable or the default value
func (a LeopardApp) GetEnv(setting string, defaultValue string) string {
	value := os.Getenv(setting)
	if value == "" {
		return defaultValue
	}
	return value
}

func (a LeopardApp) GetEnvironment() string {
	return a.GetEnv("LEOPARD_ENV", "DEVELOPMENT")
}
