package storage

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
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
