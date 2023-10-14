package config

import (
	"fmt"
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseOperations interface {
	Find(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	First(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Create(model interface{}) (tx *gorm.DB)
}

type Database struct {
	connection *gorm.DB
}

// Represents the configuration structure in the database.yml file
type envDbConfig struct {
	Development dbConfig `yaml:"development"`
	Test        dbConfig `yaml:"test"`
	Production  dbConfig `yaml:"production"`
}

// Represents the configuration for a specific database environment
type dbConfig struct {
	Dialect  string `yaml:"dialect"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Pool     int    `yaml:"pool"`
	SslMode  string `yaml:"ssl_mode"`
}

// Connects to the database based on the specified environment
func ConnectDB() DatabaseOperations {
	// Database credentials
	dbConfig := obtainDbConfig("development")

	// Create the connection string using string interpolation
	connStr := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s sslmode=%s",
		dbConfig.User,
		dbConfig.Database,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.SslMode,
	)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database")
	}

	dbs := &Database{connection: db}

	return dbs
}

// Reads the database configuration from the database.yml file
func obtainDbConfig(environment string) dbConfig {
	// Open the YAML file
	file, err := os.Open("./database.yml")
	if err != nil {
		log.Fatalf("Error opening database config file: %v", err)
	}
	defer file.Close()

	// Read the entire content of the YAML file into memory
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading database config file: %v", err)
	}

	// Define a struct to hold the parsed YAML data
	envDbConfig := envDbConfig{}

	// Unmarshal the YAML data into the struct
	err = yaml.Unmarshal([]byte(data), &envDbConfig)
	if err != nil {
		log.Fatalf("Error unmarshalling database config file: %v", err)
	}

	// Determine the appropriate database configuration based on the environment
	switch environment {
	case "development":
		return envDbConfig.Development
	case "test":
		return envDbConfig.Test
	case "production":
		return envDbConfig.Production
	default:
		// If an invalid environment is provided, panic with an error message
		panic(fmt.Errorf("invalid environment: %v", environment))
	}
}

func (db *Database) Find(dest interface{}, conds ...interface{}) (tx *gorm.DB) {
	return db.connection.Find(dest, conds...)
}

func (db *Database) First(dest interface{}, conds ...interface{}) (tx *gorm.DB) {
	return db.connection.First(dest, conds...)
}

func (db *Database) Create(value interface{}) (tx *gorm.DB) {
	return db.connection.Create(value)
}
