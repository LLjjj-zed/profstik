package service

import (
	"context"
	"github.com/132982317/profstik/kitex_gen/message"
	"github.com/132982317/profstik/pkg/utils/zap"
	Zap "go.uber.org/zap"
)

type MessageServiceImpl struct{}

func (m *MessageServiceImpl) MessageAction(ctx context.Context, req *message.MessageActionRequest) (resp *message.MessageActionResponse, err error) {
	logger := zap.LoggerPool.Get().(*Zap.SugaredLogger)
	defer zap.LoggerPool.Put(logger)
	return nil, nil
}

func (m *MessageServiceImpl) MessageChat(ctx context.Context, req *message.MessageChatRequest) (resp *message.MessageChatResponse, err error) {
	logger := zap.LoggerPool.Get().(*Zap.SugaredLogger)
	defer zap.LoggerPool.Put(logger)
	return nil, nil
}
