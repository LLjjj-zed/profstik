package service

import (
	jwt "github.com/132982317/profstik/middleware"
	"github.com/132982317/profstik/pkg/utils/rabbitmq"
	"github.com/132982317/profstik/pkg/utils/viper"
	"github.com/132982317/profstik/pkg/utils/zap"
)

var (
	Jwt        *jwt.JWT
	logger     = zap.InitLogger()
	config     = viper.Init("rabbitmq")
	autoAck    = config.Viper.GetBool("consumer.favorite.autoAck")
	FavoriteMq = rabbitmq.NewRabbitMQSimple("favorite", autoAck)
	err        error
)

func Init(signingKey string) {
	Jwt = jwt.NewJWT([]byte(signingKey))
	//GoCron()
	go consume()
}
