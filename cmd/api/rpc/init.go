package rpc

import "github.com/132982317/profstik/pkg/utils/viper"

func init() {
	// user rpc
	userConfig := viper.Init("user")
	initUserRpc(&userConfig)

	//video rpc
	videoConfig := viper.Init("video")
	initVideoRpc(&videoConfig)

	//comment rpc
	commentConfig := viper.Init("comment")
	initCommentRpc(&commentConfig)

	//favorite rpc
	favoriteConfig := viper.Init("favorite")
	initFavoriteRpc(&favoriteConfig)

	//relation rpc
	relationConfig := viper.Init("relation")
	initRelationRpc(&relationConfig)

	//message rpc
	messageConfig := viper.Init("message")
	initMessageRpc(&messageConfig)
}
