package rabbitmq

import (
	"fmt"
	"github.com/132982317/profstik/pkg/utils/viper"
	z "github.com/132982317/profstik/pkg/utils/zap"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

var (
	config = viper.Init("rabbitmq")
	logger *zap.SugaredLogger
	conn   *amqp.Connection
	err    error
	MqUrl  = fmt.Sprintf("amqp://%s:%s@%s:%d/%v",
		config.Viper.GetString("server.username"),
		config.Viper.GetString("server.password"),
		config.Viper.GetString("server.host"),
		config.Viper.GetInt("server.port"),
		config.Viper.GetString("server.vhost"),
	)
)

func init() {
	logger = z.InitLogger()
}

func failOnError(err error, msg string) {
	if err != nil {
		logger.Errorf("%s: %s", msg, err.Error())
		log.Fatal(fmt.Sprintf("%s: %s", msg, err.Error()))
	}
}
