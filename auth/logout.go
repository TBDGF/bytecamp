package auth

import (
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Logout(g *gin.RouterGroup) {
	g.Handle("POST", "/auth/logout", func(c *gin.Context) {
		// 删除 cookie
		c.SetCookie("camp-session", "", -1, "/",
			"127.0.0.1", false, true)
		var response types.LogoutResponse
		response.Code = types.OK
		c.JSON(http.StatusOK, response)
	})
}
