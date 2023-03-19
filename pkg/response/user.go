package response

import (
	"github.com/132982317/profstik/kitex_gen/user"
	"github.com/132982317/profstik/pkg/errno"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

type infoResponse struct {
	common
	User *user.User `json:"user"`
}

// UserInfoOK 返回正确信息
func UserInfoOK(c *app.RequestContext, err error, user *user.User) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, infoResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
		User: user,
	})
}

// UserInfoErr 返回错误信息
func UserInfoErr(c *app.RequestContext, err error) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, infoResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
	})
}

type LoginResponse struct {
	common
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
}

// LoginOK 返回正确信息
func LoginOK(c *app.RequestContext, err error, userid int64, token string) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, LoginResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
		UserID: userid,
		Token:  token,
	})
}

// LoginErr 返回错误信息
func LoginErr(c *app.RequestContext, err error) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, LoginResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
	})
}

type RegisterResponse struct {
	common
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
}

// RegisterOK 返回正确信息
func RegisterOK(c *app.RequestContext, err error, userid int64, token string) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, RegisterResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
		UserID: userid,
		Token:  token,
	})
}

// RegisterErr 返回错误信息
func RegisterErr(c *app.RequestContext, err error) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, RegisterResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
	})
}
