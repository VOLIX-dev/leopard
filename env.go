package leopard

import "os"

// GetEnv get a enviroment variable or the default value
func (l LeopardApp) GetEnv(setting string, defaultValue string) string {
	value := os.Getenv(setting)
	if value == "" {
		return defaultValue
	}
	return value
}

func (l LeopardApp) GetEnvironment() string {
	return l.GetEnv("LEOPARD_ENV", "DEVELOPMENT")
}
