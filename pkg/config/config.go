package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config represents application configuration
type Config struct {
	v *viper.Viper
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host         string
	Port         int
	GRPCPort     int
	ReadTimeout  int
	WriteTimeout int
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

// KafkaConfig holds Kafka configuration
type KafkaConfig struct {
	Brokers       []string
	ConsumerGroup string
	Topics        map[string]string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// OryConfig holds Ory configuration
type OryConfig struct {
	HydraAdminURL  string
	HydraPublicURL string
	KratosAdminURL string
	KratosPublicURL string
}

// ElasticsearchConfig holds Elasticsearch configuration
type ElasticsearchConfig struct {
	Addresses []string
	Username  string
	Password  string
}

// MinIOConfig holds MinIO configuration
type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
	Bucket    string
}

// New creates a new configuration instance
func New(serviceName string) (*Config, error) {
	v := viper.New()

	// Set configuration file name and paths
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("/etc/nexusflow")

	// Environment variables
	v.SetEnvPrefix(strings.ToUpper(serviceName))
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Set defaults
	setDefaults(v)

	// Read configuration file (optional)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found, will use environment variables and defaults
	}

	return &Config{v: v}, nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.grpc_port", 9090)
	v.SetDefault("server.read_timeout", 30)
	v.SetDefault("server.write_timeout", 30)

	// Database defaults
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "nexusflow")
	v.SetDefault("database.password", "nexusflow")
	v.SetDefault("database.database", "nexusflow")
	v.SetDefault("database.ssl_mode", "disable")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 5)
	v.SetDefault("database.conn_max_lifetime", 300)

	// Kafka defaults
	v.SetDefault("kafka.brokers", []string{"localhost:9092"})
	v.SetDefault("kafka.consumer_group", "nexusflow")

	// Redis defaults
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.db", 0)

	// Elasticsearch defaults
	v.SetDefault("elasticsearch.addresses", []string{"http://localhost:9200"})

	// MinIO defaults
	v.SetDefault("minio.endpoint", "localhost:9000")
	v.SetDefault("minio.use_ssl", false)
	v.SetDefault("minio.bucket", "nexusflow")
}

// GetServer returns server configuration
func (c *Config) GetServer() ServerConfig {
	return ServerConfig{
		Host:         c.v.GetString("server.host"),
		Port:         c.v.GetInt("server.port"),
		GRPCPort:     c.v.GetInt("server.grpc_port"),
		ReadTimeout:  c.v.GetInt("server.read_timeout"),
		WriteTimeout: c.v.GetInt("server.write_timeout"),
	}
}

// GetDatabase returns database configuration
func (c *Config) GetDatabase() DatabaseConfig {
	return DatabaseConfig{
		Host:            c.v.GetString("database.host"),
		Port:            c.v.GetInt("database.port"),
		User:            c.v.GetString("database.user"),
		Password:        c.v.GetString("database.password"),
		Database:        c.v.GetString("database.database"),
		SSLMode:         c.v.GetString("database.ssl_mode"),
		MaxOpenConns:    c.v.GetInt("database.max_open_conns"),
		MaxIdleConns:    c.v.GetInt("database.max_idle_conns"),
		ConnMaxLifetime: c.v.GetInt("database.conn_max_lifetime"),
	}
}

// GetKafka returns Kafka configuration
func (c *Config) GetKafka() KafkaConfig {
	return KafkaConfig{
		Brokers:       c.v.GetStringSlice("kafka.brokers"),
		ConsumerGroup: c.v.GetString("kafka.consumer_group"),
		Topics:        c.v.GetStringMapString("kafka.topics"),
	}
}

// GetRedis returns Redis configuration
func (c *Config) GetRedis() RedisConfig {
	return RedisConfig{
		Host:     c.v.GetString("redis.host"),
		Port:     c.v.GetInt("redis.port"),
		Password: c.v.GetString("redis.password"),
		DB:       c.v.GetInt("redis.db"),
	}
}

// GetOry returns Ory configuration
func (c *Config) GetOry() OryConfig {
	return OryConfig{
		HydraAdminURL:   c.v.GetString("ory.hydra_admin_url"),
		HydraPublicURL:  c.v.GetString("ory.hydra_public_url"),
		KratosAdminURL:  c.v.GetString("ory.kratos_admin_url"),
		KratosPublicURL: c.v.GetString("ory.kratos_public_url"),
	}
}

// GetElasticsearch returns Elasticsearch configuration
func (c *Config) GetElasticsearch() ElasticsearchConfig {
	return ElasticsearchConfig{
		Addresses: c.v.GetStringSlice("elasticsearch.addresses"),
		Username:  c.v.GetString("elasticsearch.username"),
		Password:  c.v.GetString("elasticsearch.password"),
	}
}

// GetMinIO returns MinIO configuration
func (c *Config) GetMinIO() MinIOConfig {
	return MinIOConfig{
		Endpoint:  c.v.GetString("minio.endpoint"),
		AccessKey: c.v.GetString("minio.access_key"),
		SecretKey: c.v.GetString("minio.secret_key"),
		UseSSL:    c.v.GetBool("minio.use_ssl"),
		Bucket:    c.v.GetString("minio.bucket"),
	}
}

// Get returns a configuration value
func (c *Config) Get(key string) interface{} {
	return c.v.Get(key)
}

// GetString returns a string configuration value
func (c *Config) GetString(key string) string {
	return c.v.GetString(key)
}

// GetInt returns an int configuration value
func (c *Config) GetInt(key string) int {
	return c.v.GetInt(key)
}

// GetBool returns a bool configuration value
func (c *Config) GetBool(key string) bool {
	return c.v.GetBool(key)
}
