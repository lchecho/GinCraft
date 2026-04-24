package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config 全局配置实例
var Config AppConfig

// AppConfig 应用配置
type AppConfig struct {
	App struct {
		Name         string `mapstructure:"name"`
		Mode         string `mapstructure:"mode"`
		Port         int    `mapstructure:"port"`
		ReadTimeout  int    `mapstructure:"read_timeout"`
		WriteTimeout int    `mapstructure:"write_timeout"`
		APIKey       string `mapstructure:"api_key"`
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
		Host         string `mapstructure:"host"`
		Port         int    `mapstructure:"port"`
		Password     string `mapstructure:"password"`
		DB           int    `mapstructure:"db"`
		PoolSize     int    `mapstructure:"pool_size"`
		MinIdleConns int    `mapstructure:"min_idle_conns"`
		MaxRetries   int    `mapstructure:"max_retries"`
		DialTimeout  int    `mapstructure:"dial_timeout"`
		ReadTimeout  int    `mapstructure:"read_timeout"`
		WriteTimeout int    `mapstructure:"write_timeout"`
	} `mapstructure:"redis"`

	CORS struct {
		AllowOrigins     []string `mapstructure:"allow_origins"` // 空或含 "*" 表示允许所有来源
		AllowMethods     []string `mapstructure:"allow_methods"`
		AllowHeaders     []string `mapstructure:"allow_headers"`
		AllowCredentials bool     `mapstructure:"allow_credentials"`
		MaxAge           int      `mapstructure:"max_age"`
	} `mapstructure:"cors"`
}

// LoadConfig 加载配置
func LoadConfig(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("read config: %w", err)
	}
	if err := viper.Unmarshal(&Config); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}
	if err := validate(); err != nil {
		return err
	}

	logDir := filepath.Dir(Config.Log.Filename)
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return fmt.Errorf("mkdir log dir: %w", err)
	}
	return nil
}

func setDefaults() {
	viper.SetDefault("app.mode", "debug")
	viper.SetDefault("app.port", 8080)
	viper.SetDefault("app.read_timeout", 60)
	viper.SetDefault("app.write_timeout", 60)

	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.filename", "logs/app.log")
	viper.SetDefault("log.max_size", 100)
	viper.SetDefault("log.max_age", 30)
	viper.SetDefault("log.max_backups", 10)

	viper.SetDefault("mysql.max_idle_conns", 10)
	viper.SetDefault("mysql.max_open_conns", 100)
	viper.SetDefault("mysql.conn_max_lifetime", 3600)

	viper.SetDefault("redis.pool_size", 10)
	viper.SetDefault("redis.min_idle_conns", 5)
	viper.SetDefault("redis.max_retries", 3)
	viper.SetDefault("redis.dial_timeout", 5)
	viper.SetDefault("redis.read_timeout", 3)
	viper.SetDefault("redis.write_timeout", 3)

	viper.SetDefault("cors.allow_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("cors.allow_headers", []string{
		"Content-Type", "Authorization", "X-Requested-With", "X-API-Key", "X-Trace-ID",
	})
	viper.SetDefault("cors.max_age", 600)
}

func validate() error {
	if Config.App.Port <= 0 {
		return fmt.Errorf("config: app.port must be > 0")
	}
	if Config.MySQL.Host == "" || Config.MySQL.Database == "" {
		return fmt.Errorf("config: mysql.host and mysql.database are required")
	}
	return nil
}
