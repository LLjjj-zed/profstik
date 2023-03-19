package response

import (
	"github.com/132982317/profstik/kitex_gen/relation"
	"github.com/132982317/profstik/kitex_gen/user"
	"github.com/132982317/profstik/pkg/errno"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

type RelationResponse struct {
	common
}

// RelationOk 返回正确信息
func RelationOk(c *app.RequestContext, err error) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, &RelationResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
	})
}

// RelationErr  返回错误信息
func RelationErr(c *app.RequestContext, err error) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, &RelationResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
	})
}

type FollowerListResponse struct {
	common
	UserList []*user.User `json:"user_list"`
}

// FollowerListOk 返回正确信息
func FollowerListOk(c *app.RequestContext, err error, UserList []*user.User) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, &FollowerListResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
		UserList: UserList,
	})
}

// FollowerListErr  返回错误信息
func FollowerListErr(c *app.RequestContext, err error) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, &FollowerListResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
	})
}

type FollowListResponse struct {
	common
	UserList []*user.User `json:"user_list"`
}

// FollowListOk 返回正确信息
func FollowListOk(c *app.RequestContext, err error, UserList []*user.User) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, &FollowListResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
		UserList: UserList,
	})
}

// FollowListErr  返回错误信息
func FollowListErr(c *app.RequestContext, err error) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, &FollowListResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
	})
}

type FriendListResponse struct {
	common
	UserList []*relation.FriendUser `json:"user_list"`
}

// FriendListOk 返回正确信息
func FriendListOk(c *app.RequestContext, err error, UserList []*relation.FriendUser) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, &FriendListResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
		UserList: UserList,
	})
}

// FriendListErr  返回错误信息
func FriendListErr(c *app.RequestContext, err error) {
	Err := errno.ConvertErr(err)
	c.JSON(http.StatusOK, &FriendListResponse{
		common: common{
			StatusCode: Err.ErrCode,
			StatusMsg:  Err.ErrMsg,
		},
	})
}
