package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/joho/godotenv"
)

type Secrets interface {
	GetConnParams() (*postgresParams, error)
}

type secretsManager struct {
	client *secretsmanager.SecretsManager
	arn    string
}

func NewSecretsManager(awsSession *session.Session, arn string) *secretsManager {
	client := secretsmanager.New(awsSession)
	return &secretsManager{client: client, arn: arn}
}

func (s *secretsManager) GetConnParams() (*postgresParams, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(s.arn),
	}
	value, err := s.client.GetSecretValue(input)
	if err != nil {
		return nil, err
	}
	params := &postgresParams{}
	if err := json.Unmarshal([]byte(*value.SecretString), params); err != nil {
		return nil, err
	}
	params.ssl = "require"
	return params, nil
}

type envLoader struct{}

func NewEnvLoader() (*envLoader, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	return &envLoader{}, nil
}

func (l *envLoader) GetConnParams() (*postgresParams, error) {
	host, err := envOrError("DB_HOST")
	if err != nil {
		return nil, err
	}
	port, err := envOrError("DB_PORT")
	if err != nil {
		return nil, err
	}
	name, err := envOrError("DB_NAME")
	if err != nil {
		return nil, err
	}
	user, err := envOrError("DB_USER")
	if err != nil {
		return nil, err
	}
	pass, err := envOrError("DB_PASS")
	if err != nil {
		return nil, err
	}
	return &postgresParams{
		DBHost: host,
		DBPort: port,
		DBName: name,
		DBUser: user,
		DBPass: pass,
		ssl:    "disable",
	}, nil
}

func envOrError(key string) (string, error) {
	value := os.Getenv(key)
	if len(value) < 1 {
		return "", errors.New(fmt.Sprintf("Environment variable '%s' must be specified", key))
	}
	return value, nil
}
