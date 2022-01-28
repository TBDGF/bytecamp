package auth

import (
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
)

func Logout(g *gin.Engine, Db *sqlx.DB) {
	g.Handle("POST", "/auth/logout", func(c *gin.Context) {
		// 删除 cookie
		c.SetCookie("camp-session", "", -1, "/auth/whoami",
			"127.0.0.1", false, true)
		var response types.LogoutResponse
		response.Code = 0
		c.JSON(http.StatusOK, response)
	})
}
