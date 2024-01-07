package storage

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/joho/godotenv"
)

type Secrets interface {
	GetConnParams() (*postgresConnParams, error)
}

type secretsManager struct {
	client *secretsmanager.SecretsManager
	arn    string
}

func NewSecretsManager(awsSession *session.Session, arn string) *secretsManager {
	client := secretsmanager.New(awsSession)
	return &secretsManager{client: client, arn: arn}
}

func (s *secretsManager) GetConnParams() (*postgresConnParams, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(s.arn),
	}
	value, err := s.client.GetSecretValue(input)
	if err != nil {
		return nil, err
	}
	type awsSecrets struct {
		DBHost string `json:"DB_HOST"`
		DBPort string `json:"DB_PORT"`
		DBUser string `json:"DB_USER"`
		DBPass string `json:"DB_PASS"`
		DBName string `json:"DB_NAME"`
	}
	secrets := &awsSecrets{}
	if err := json.Unmarshal([]byte(*value.SecretString), secrets); err != nil {
		return nil, err
	}
	if len(secrets.DBHost) < 6 {
		return nil, errors.New("'DB_HOST' secret corrupted")
	}
	connParams := &postgresConnParams{
		DBHost: secrets.DBHost[:len(secrets.DBHost)-5],
		DBPort: secrets.DBPort,
		DBName: secrets.DBName,
		DBUser: secrets.DBUser,
		DBPass: secrets.DBPass,
		SSL:    "require",
	}
	return connParams, nil
}

type envLoader struct{}

func NewEnvLoader() (*envLoader, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	return &envLoader{}, nil
}

func (l *envLoader) GetConnParams() (*postgresConnParams, error) {
	dbHost, err := getEnv("DB_HOST")
	if err != nil {
		return nil, err
	}
	dbPort, err := getEnv("DB_PORT")
	if err != nil {
		return nil, err
	}
	dbName, err := getEnv("DB_NAME")
	if err != nil {
		return nil, err
	}
	dbUser, err := getEnv("DB_USER")
	if err != nil {
		return nil, err
	}
	dbPass, err := getEnv("DB_PASS")
	if err != nil {
		return nil, err
	}
	return &postgresConnParams{
		DBHost: dbHost,
		DBPort: dbPort,
		DBName: dbName,
		DBUser: dbUser,
		DBPass: dbPass,
		SSL:    "disable",
	}, nil
}

func getEnv(key string) (string, error) {
	value := os.Getenv(key)
	if len(value) < 1 {
		return "", errors.New("Environment variable '" + key + "' must be specified")
	}
	return value, nil
}
