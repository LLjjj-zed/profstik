package webHandler

import (
	"context"
	"github.com/132982317/profstik/cmd/api/rpc"
	kitex "github.com/132982317/profstik/kitex_gen/user"
	"github.com/132982317/profstik/pkg/errno"
	"github.com/132982317/profstik/pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
)

// Register 注册
func Register(ctx context.Context, c *app.RequestContext) {
	var registerVar RegisterParam
	// 绑定请求体中的JSON数据到RegisterParam结构体
	if err := c.Bind(&registerVar); err != nil {
		response.RegisterErr(c, errno.ConvertErr(err))
		return
	}
	// 校验请求参数
	if len(registerVar.Username) == 0 || len(registerVar.Password) == 0 || len(registerVar.Username) > 32 || len(registerVar.Password) > 32 || len(registerVar.Password) < 5 {
		response.RegisterErr(c, errno.ParamErr)
		return
	}
	//调用rpc
	register, err := rpc.Register(ctx, &kitex.UserRegisterRequest{
		Username: registerVar.Username,
		Password: registerVar.Password,
	})
	if err != nil {
		response.RegisterErr(c, errno.RpcConnectErr)
		return
	}
	response.RegisterOK(c, errno.Success, register.UserId, register.Token)
}

// Login 登录
func Login(ctx context.Context, c *app.RequestContext) {
	var loginVar LoginParam
	// 绑定请求体中的JSON数据到RegisterParam结构体
	if err := c.Bind(&loginVar); err != nil {
		response.LoginErr(c, errno.ConvertErr(err))
		return
	}
	// 校验请求参数
	if len(loginVar.Username) == 0 || len(loginVar.Password) == 0 {
		response.LoginErr(c, errno.ParamErr)
		return
	}
	login, err := rpc.Login(ctx, &kitex.UserLoginRequest{
		Username: loginVar.Username,
		Password: loginVar.Password,
	})
	if err != nil {
		response.LoginErr(c, errno.RpcConnectErr)
		return
	}
	response.LoginOK(c, errno.Success, login.UserId, login.Token)
}

func UserInfo(ctx context.Context, c *app.RequestContext) {
	var userinfoVar UserInfoParam
	// 绑定请求体中的JSON数据到RegisterParam结构体
	if err := c.Bind(&userinfoVar); err != nil {
		response.UserInfoErr(c, errno.ConvertErr(err))
		return
	}
	if len(userinfoVar.Token) == 0 {
		response.UserInfoErr(c, errno.AuthorizationFailedErr)
	}
	info, err := rpc.UserInfo(ctx, &kitex.UserInfoRequest{
		UserId: userinfoVar.UserID,
		Token:  userinfoVar.Token,
	})
	if err != nil {
		response.UserInfoErr(c, errno.RpcConnectErr)
		return
	}
	response.UserInfoOK(c, errno.Success, info.User)
}
