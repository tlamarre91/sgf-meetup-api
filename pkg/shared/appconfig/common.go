package appconfig

import (
	"context"
	"github.com/google/wire"
	"sgf-meetup-api/pkg/shared/db"
	"sgf-meetup-api/pkg/shared/logging"
)

const (
	LogLevelKey                   = "LOG_LEVEL"
	LogTypeKey                    = "LOG_TYPE"
	SentryDSNKey                  = "SENTRY_DSN"
	appEnvKey                     = "APP_ENV"
	SSMPathKey                    = "SSM_PATH"
	AWSRegionKey                  = "AWS_REGION"
	AWSAccessKeyKey               = "AWS_ACCESS_KEY"
	AWSSecretAccessKeyKey         = "AWS_SECRET_ACCESS_KEY"
	AWSSessionTokenKey            = "AWS_SESSION_TOKEN"
	AWSProfileKey                 = "AWS_PROFILE"
	AWSConfigFileKey              = "AWS_CONFIG_FILE"
	AWSSharedCredentialsFileKey   = "AWS_SHARED_CREDENTIALS_FILE"
	DynamoDBEndpointKey           = "DYNAMODB_ENDPOINT"
	DynamoDBAWSRegionKey          = "DYNAMODB_AWS_REGION"
	DynamoDBAWSAccessKeyKey       = "DYNAMODB_AWS_ACCESS_KEY"
	DynamoDBAWSSecretAccessKeyKey = "DYNAMODB_AWS_SECRET_ACCESS_KEY"
)

var CommonKeys = []string{
	LogLevelKey,
	LogTypeKey,
	SentryDSNKey,
	appEnvKey,
	SSMPathKey,
	AWSRegionKey,
	AWSAccessKeyKey,
	AWSSecretAccessKeyKey,
	AWSSessionTokenKey,
	AWSProfileKey,
	AWSConfigFileKey,
	AWSSharedCredentialsFileKey,
	DynamoDBEndpointKey,
	DynamoDBAWSRegionKey,
	DynamoDBAWSAccessKeyKey,
	DynamoDBAWSSecretAccessKeyKey,
}

type Aws struct {
	AwsRegion          string `mapstructure:"aws_region"`
	AwsAccessKey       string `mapstructure:"aws_access_key"`
	AwsSecretAccessKey string `mapstructure:"aws_secret_access_key"`
}

type DynamoDB struct {
	DynamoDbEndpoint   string `mapstructure:"dynamodb_endpoint"`
	AwsRegion          string `mapstructure:"dynamodb_aws_region"`
	AwsAccessKey       string `mapstructure:"dynamodb_aws_access_key"`
	AwsSecretAccessKey string `mapstructure:"dynamodb_aws_secret_access_key"`
}

type Common struct {
	Logging   logging.Config `mapstructure:",squash"`
	SentryDSN string         `mapstructure:"sentry_dsn"`
	AppEnv    string         `mapstructure:"app_env"`
	Aws       Aws            `mapstructure:",squash"`
	DynamoDB  db.Config      `mapstructure:",squash"`
}

var ConfigProviders = wire.NewSet(
	wire.Bind(new(AwsConfigManager), new(*AwsConfigManagerImpl)),
	NewAwsConfigManager,
	AwsConfigProvider,
	wire.FieldsOf(new(Common), "Logging", "Aws", "DynamoDB"),
)

func NewCommonConfig(ctx context.Context) (*Common, error) {
	var config Common

	err := NewParser().
		WithCommonConfig().
		WithEnvFile(".", ".env").
		WithEnvVars().
		Parse(ctx, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
