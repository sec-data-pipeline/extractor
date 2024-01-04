package storage

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type Secrets struct {
	DBHost string `json:"DB_HOST"`
	DBPort string `json:"DB_PORT"`
	DBUser string `json:"DB_USER"`
	DBPass string `json:"DB_PASS"`
	DBName string `json:"DB_NAME"`
}

func GetAWSConnParams(awsSession *session.Session, secretsARN string) (*DBConnParams, error) {
	client := secretsmanager.New(awsSession)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretsARN),
	}
	result, err := client.GetSecretValue(input)
	if err != nil {
		return nil, err
	}
	secrets := &Secrets{}
	if err := json.Unmarshal([]byte(*result.SecretString), secrets); err != nil {
		return nil, err
	}
	if len(secrets.DBHost) < 6 {
		return nil, errors.New("'DB_HOST' secret corrupted")
	}
	connParams := &DBConnParams{
		DBHost: secrets.DBHost[:len(secrets.DBHost)-5],
		DBPort: secrets.DBPort,
		DBName: secrets.DBName,
		DBUser: secrets.DBUser,
		DBPass: secrets.DBPass,
		SSL:    "require",
	}
	return connParams, nil
}
