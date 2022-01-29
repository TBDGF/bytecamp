package auth

import (
	"bytedance/config"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Logout(c *gin.Context) {
	// 删除 cookie
	c.SetCookie("camp-session", "", -1, "/",
		config.Host, false, true)
	var response types.LogoutResponse
	response.Code = types.OK
	c.JSON(http.StatusOK, response)
}
