package main

import (
	"context"
	"github.com/132982317/profstik/kitex_gen/message"
)

type MessageServiceImpl struct{}

func (m *MessageServiceImpl) MessageAction(ctx context.Context, req *message.MessageActionRequest) (*message.MessageActionResponse, error) {
	return nil, nil
}

func (m *MessageServiceImpl) MessageChat(ctx context.Context, req *message.MessageChatRequest) (*message.MessageChatResponse, error) {
	return nil, nil
}
