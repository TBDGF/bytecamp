package auth

import (
	"bytedance/db"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Logout(c *gin.Context) {
	var response types.LogoutResponse

	// 判断用户是否登录
	_, err := c.Cookie("camp-session")
	if err != nil {
		response.Code = types.LoginRequired
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// 删除 cookie
	c.SetCookie("camp-session", "", -1, "/",
		db.Host, false, true)
	response.Code = types.OK
	c.JSON(http.StatusOK, response)
}
