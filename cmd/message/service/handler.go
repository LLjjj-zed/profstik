package service

import (
	"context"
	"github.com/132982317/profstik/dao/mysql"
	"github.com/132982317/profstik/dao/redis"
	"github.com/132982317/profstik/kitex_gen/message"
	"github.com/132982317/profstik/pkg/errno"
	tool "github.com/132982317/profstik/pkg/utils/crypt"
	"github.com/132982317/profstik/pkg/utils/zap"
	Zap "go.uber.org/zap"
)

type MessageServiceImpl struct{}

func (m *MessageServiceImpl) MessageAction(ctx context.Context, req *message.MessageActionRequest) (resp *message.MessageActionResponse, err error) {
	logger := zap.LoggerPool.Get().(*Zap.SugaredLogger)
	defer zap.LoggerPool.Put(logger)

	// 解析token,获取用户id
	claims, err := Jwt.ParseToken(req.Token)
	if err != nil {
		logger.Errorln(err.Error())
		resp = &message.MessageActionResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	userID := claims.Id

	toUserID, actionType := req.ToUserId, req.ActionType

	if userID == toUserID {
		logger.Errorln("不能给自己发送消息")
		resp = &message.MessageActionResponse{
			StatusCode: errno.ServiceErrCode,
			StatusMsg:  "消息发送失败：不能给自己发送消息",
		}
		return
	}

	relation, err := mysql.GetRelationByUserIDs(ctx, userID, toUserID)
	if relation == nil {
		logger.Errorf("消息发送失败：非朋友关系，无法发送")
		resp = &message.MessageActionResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}

	rsaContent, err := tool.RsaEncrypt([]byte(req.Content), publicKey)
	if err != nil {
		logger.Errorf("rsa encrypt error: %v\n", err.Error())
		resp = &message.MessageActionResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}

	messages := make([]*mysql.Message, 0)
	messages = append(messages, &mysql.Message{
		FromUserID: uint(userID),
		ToUserID:   uint(toUserID),
		Content:    string(tool.Base64Encode(rsaContent)),
	})
	if actionType == 1 {
		err := mysql.CreateMessagesByList(ctx, messages)
		if err != nil {
			logger.Errorln(err.Error())
			resp = &message.MessageActionResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
	} else {
		logger.Errorf("action_type 非法：%v", actionType)
		resp = &message.MessageActionResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	resp = &message.MessageActionResponse{
		StatusCode: errno.SuccessCode,
	}
	return resp, nil
}

func (m *MessageServiceImpl) MessageChat(ctx context.Context, req *message.MessageChatRequest) (resp *message.MessageChatResponse, err error) {
	logger := zap.LoggerPool.Get().(*Zap.SugaredLogger)
	defer zap.LoggerPool.Put(logger)
	// 解析token,获取用户id
	claims, err := Jwt.ParseToken(req.Token)
	if err != nil {
		logger.Errorln(err.Error())
		resp = &message.MessageChatResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	userID := claims.Id

	// 从redis中获取message时间戳
	lastTimestamp, err := redis.GetMessageTimestamp(ctx, req.Token, req.ToUserId)
	if err != nil {
		logger.Errorln(err.Error())
		resp = &message.MessageChatResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}

	var results []*mysql.Message
	if lastTimestamp == -1 {
		results, err = mysql.GetMessagesByUserIDs(ctx, userID, req.ToUserId, int64(lastTimestamp))
		lastTimestamp = 0
	} else {
		results, err = mysql.GetMessagesByUserToUser(ctx, req.ToUserId, userID, int64(lastTimestamp))
	}

	if err != nil {
		logger.Errorln(err.Error())
		resp = &message.MessageChatResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	messages := make([]*message.Message, 0)
	for _, r := range results {
		decContent, err := tool.Base64Decode([]byte(r.Content))
		if err != nil {
			logger.Errorf("Base64Decode error: %v\n", err.Error())
			resp = &message.MessageChatResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
		decContent, err = tool.RsaDecrypt(decContent, privateKey)
		if err != nil {
			logger.Errorf("rsa decrypt error: %v\n", err.Error())
			resp = &message.MessageChatResponse{
				StatusCode: errno.ServiceErrCode,
			}
			return
		}
		messages = append(messages, &message.Message{
			Id:         int64(r.ID),
			FromUserId: int64(r.FromUserID),
			ToUserId:   int64(r.ToUserID),
			Content:    string(decContent),
			CreateTime: r.CreatedAt.UnixMilli(),
		})
	}

	resp = &message.MessageChatResponse{
		StatusCode:  errno.SuccessCode,
		MessageList: messages,
	}

	// 更新时间redis里的时间戳
	if len(messages) > 0 {
		Message := messages[len(messages)-1]
		lastTimestamp = int(Message.CreateTime)
	}

	if err = redis.SetMessageTimestamp(ctx, req.Token, req.ToUserId, lastTimestamp); err != nil {
		logger.Errorln(err.Error())
		resp = &message.MessageChatResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}

	return resp, nil
}
