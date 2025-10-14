package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var Config AppConfig

type AppConfig struct {
	App struct {
		Name         string `mapstructure:"name"`
		Mode         string `mapstructure:"mode"`
		Port         int    `mapstructure:"port"`
		ReadTimeout  int    `mapstructure:"read_timeout"`
		WriteTimeout int    `mapstructure:"write_timeout"`
	} `mapstructure:"app"`

	Log struct {
		Level      string `mapstructure:"level"`
		Filename   string `mapstructure:"filename"`
		MaxSize    int    `mapstructure:"max_size"`
		MaxAge     int    `mapstructure:"max_age"`
		MaxBackups int    `mapstructure:"max_backups"`
		Compress   bool   `mapstructure:"compress"`
	} `mapstructure:"log"`

	MySQL struct {
		Host            string `mapstructure:"host"`
		Port            int    `mapstructure:"port"`
		Username        string `mapstructure:"username"`
		Password        string `mapstructure:"password"`
		Database        string `mapstructure:"database"`
		MaxIdleConns    int    `mapstructure:"max_idle_conns"`
		MaxOpenConns    int    `mapstructure:"max_open_conns"`
		ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	} `mapstructure:"mysql"`

	Redis struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
		PoolSize int    `mapstructure:"pool_size"`
	} `mapstructure:"redis"`

	Trace struct {
		Enabled bool   `mapstructure:"enabled"`
		Type    string `mapstructure:"type"`
		Jaeger  struct {
			ServiceName       string  `mapstructure:"service_name"`
			CollectorEndpoint string  `mapstructure:"collector_endpoint"`
			SamplerType       string  `mapstructure:"sampler_type"`
			SamplerParam      float64 `mapstructure:"sampler_param"`
		} `mapstructure:"jaeger"`
	} `mapstructure:"trace"`
}

// LoadConfig 加载配置
func LoadConfig(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := viper.Unmarshal(&Config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 确保日志目录存在
	logDir := filepath.Dir(Config.Log.Filename)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	return nil
}
