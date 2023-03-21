package service

import (
	"context"
	"fmt"
	"github.com/132982317/profstik/dao/mysql"
	"github.com/132982317/profstik/kitex_gen/user"
	jwt "github.com/132982317/profstik/middleware"
	"github.com/132982317/profstik/pkg/errno"
	"github.com/132982317/profstik/pkg/utils/aes"
	"github.com/132982317/profstik/pkg/utils/zap"
	Zap "go.uber.org/zap"
	"math/rand"
	"time"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

func (u *UserServiceImpl) Register(ctx context.Context, req *user.UserRegisterRequest) (resp *user.UserRegisterResponse, err error) {
	logger := zap.LoggerPool.Get().(*Zap.SugaredLogger)
	defer zap.LoggerPool.Put(logger)
	if _, err = mysql.GetUserByUserName(ctx, req.Username); err != nil {
		logger.Errorln(err.Error())
		resp = &user.UserRegisterResponse{
			StatusCode: errno.UserAlreadyExistErrCode,
		}
		return
	}
	pwd, err := aes.EnPwdCode(req.Password)
	if err != nil {
		logger.Errorln(err.Error())
		resp = &user.UserRegisterResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	errno.Dprintf("[Register] pwd:%s", pwd)
	register := &mysql.User{
		UserName: req.Username,
		Password: pwd,
		Avatar:   fmt.Sprintf("default%d.png", rand.Intn(10)),
	}
	if err := mysql.CreateUser(ctx, register); err != nil {
		return nil, err
	}
	claims := jwt.CustomClaims{Id: int64(register.ID)}
	claims.ExpiresAt = time.Now().Add(time.Minute * 5).Unix()
	token, err := Jwt.CreateToken(claims)
	if err != nil {
		logger.Errorln(err.Error())
		resp = &user.UserRegisterResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	errno.Dprintf("[Register] token:%s", token)
	resp = &user.UserRegisterResponse{
		StatusCode: errno.SuccessCode,
		UserId:     int64(register.ID),
		Token:      token,
	}
	return resp, nil
}

func (u *UserServiceImpl) Login(ctx context.Context, req *user.UserLoginRequest) (resp *user.UserLoginResponse, err error) {
	logger := zap.LoggerPool.Get().(*Zap.SugaredLogger)
	defer zap.LoggerPool.Put(logger)

	// 根据用户名获取密码
	usr, err := mysql.GetUserByUserName(ctx, req.Username)
	if err != nil {
		logger.Errorln(err.Error())
		resp = &user.UserLoginResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	} else if usr == nil {
		resp = &user.UserLoginResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	errno.Dprintf("[Login] user:%+v", usr)
	// 比较数据库中的密码和请求的密码
	pwd, err := aes.DePwdCode(req.Password)
	if err != nil {
		logger.Errorln(err.Error())
		resp = &user.UserLoginResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	if string(pwd) != usr.Password {
		logger.Errorln("用户名或密码错误")
		resp = &user.UserLoginResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}
	errno.Dprintf("[Login] user:%+v", usr)
	// 密码认证通过,获取用户id并生成token
	claims := jwt.CustomClaims{
		Id: int64(usr.ID),
	}
	claims.ExpiresAt = time.Now().Add(time.Hour * 24).Unix()
	token, err := Jwt.CreateToken(claims)
	errno.Dprintf("[Login] token:%s", token)
	if err != nil {
		logger.Errorf("发生错误：%v", err.Error())
		resp = &user.UserLoginResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}

	// 返回结果
	resp = &user.UserLoginResponse{
		StatusCode: errno.SuccessCode,
		UserId:     int64(usr.ID),
		Token:      token,
	}
	return resp, nil
}

func (u *UserServiceImpl) UserInfo(ctx context.Context, req *user.UserInfoRequest) (resp *user.UserInfoResponse, err error) {
	logger := zap.LoggerPool.Get().(*Zap.SugaredLogger)
	defer zap.LoggerPool.Put(logger)

	userID := req.UserId
	// 从数据库获取user
	usr, err := mysql.GetUserByID(ctx, userID)
	if err != nil {
		logger.Errorf("发生错误：%v", err.Error())
		resp = &user.UserInfoResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	} else if usr == nil {
		logger.Errorf("该用户不存在：%v", err.Error())
		resp = &user.UserInfoResponse{
			StatusCode: errno.ServiceErrCode,
		}
		return
	}

	//返回结果
	resp = &user.UserInfoResponse{
		StatusCode: errno.ServiceErrCode,
		User: &user.User{
			Id:              int64(usr.ID),
			Name:            usr.UserName,
			FollowCount:     int64(usr.FollowingCount),
			FollowerCount:   int64(usr.FollowerCount),
			IsFollow:        userID == int64(usr.ID),
			Avatar:          "avatar",
			BackgroundImage: "backgroundImage",
			Signature:       usr.Signature,
			TotalFavorited:  int64(usr.TotalFavorited),
			WorkCount:       int64(usr.WorkCount),
			FavoriteCount:   int64(usr.FavoriteCount),
		},
	}
	return resp, nil
}
