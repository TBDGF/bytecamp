package member

import (
	"bytedance/config"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Update(c *gin.Context) {
	var request types.UpdateMemberRequest
	var response types.UpdateMemberResponse

	// -----验证操作权限 : 无权限返回 PermDenied ------ //
	// 根据 cookie 获取当前用户权限
	cookie, err := c.Cookie("camp-session")
	if err != nil {
		response.Code = types.LoginRequired // cookie 过期，用户未登录
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var usertype []types.UserType
	config.NewDB().Select(&usertype, "select usertype from userinfo where userid = ?", cookie)
	if usertype[0] != types.Admin {
		response.Code = types.PermDenied
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// ---- 验证用户昵称: 不小于 4 位，不超过 20 位 ----
	if len(request.Nickname) < 4 || len(request.Nickname) > 20 {
		response.Code = types.ParamInvalid
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, QErr := config.NewDB().Exec("update userinfo set nickname = ? where userid = ?", request.Nickname, request.UserID)
	if QErr != nil {
		// 用户不存在
		response.Code = types.UserNotExisted
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response.Code = types.OK
	c.JSON(http.StatusOK, response)
	return
}
