package member

import (
	"bytedance/config"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Delete(c *gin.Context) {
	var request types.DeleteMemberRequest
	var response types.DeleteMemberResponse

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

	// 删除成员
	ret, _ := config.NewDB().Exec(
		"delete from users where name = (select username from userinfo where userid = ?)", request.UserID)
	row, _ := ret.RowsAffected()

	if row == 0 {
		// 用户不存在
		response.Code = types.UserNotExisted
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response.Code = types.OK
	c.JSON(http.StatusOK, response)
	return
}
