package service

import (
	jwt "github.com/132982317/profstik/middleware"
	tool "github.com/132982317/profstik/pkg/utils/crypt"
	"github.com/132982317/profstik/pkg/utils/rabbitmq"
	"github.com/132982317/profstik/pkg/utils/viper"
	"github.com/132982317/profstik/pkg/utils/zap"
)

var (
	Jwt        *jwt.JWT
	logger     = zap.InitLogger()
	config     = viper.Init("rabbitmq")
	autoAck    = config.Viper.GetBool("consumer.relation.autoAck")
	RelationMq = rabbitmq.NewRabbitMQSimple("relation", autoAck)
	err        error
	privateKey string
)

func Init(signingKey string) {
	Jwt = jwt.NewJWT([]byte(signingKey))
	privateKey, _ = tool.ReadKeyFromFile(tool.PrivateKeyFilePath)
	//GoCron()
	go consume()
}
