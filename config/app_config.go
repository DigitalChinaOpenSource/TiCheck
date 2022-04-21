package config

var GlobalConfig *AppConfig

type AppConfig struct {
	WorkDir string
	Port    int
}

func (c *AppConfig) GetProbePrefix() string {
	return c.WorkDir + "probes"
}

func (c *AppConfig) GetStorePath() string {
	return c.WorkDir + "store/ticheck.db"
}
