package polly

import (
	"context"
	"github.com/Fufuhu/gopolly/util/logging"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	"go.uber.org/zap"
)

type ClientConfig struct {
}

type ClientError struct{}

func (e *ClientError) Error() string {
	return "failed to create client"
}

var client *polly.Client

// GetPollyClient Pollyのクライアントを取得する
func GetPollyClient(clientConfig *ClientConfig) (*polly.Client, error) {
	logger := logging.GetLogger()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)

	var cfg aws.Config
	var err error
	if client == nil {
		if cfg, err = config.LoadDefaultConfig(context.TODO()); err != nil {
			logger.Warn(err.Error())
			return nil, err
		}
	}
	if client = polly.NewFromConfig(cfg); client == nil {
		logger.Warn("failed to create client")
		return nil, &ClientError{}
	}
	return client, nil
}
