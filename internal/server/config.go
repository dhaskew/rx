package server

import (
	"fmt"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

const DefaultEnvPath = "./config/local.env"

type Config struct {
	path   string
	Logger *zap.Logger
	env    map[string]string
}

func (c *Config) Load(envPath string) (map[string]string, error) {
	c.path = envPath
	c.Logger.Info("Using env config file: " + c.path)
	env, err := godotenv.Read(c.path)
	c.env = env
	cnt := len(c.env)
	c.Logger.Info(fmt.Sprintf("Loaded %d environment variables", cnt))
	return c.env, err
}
