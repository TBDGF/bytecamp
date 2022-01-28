package auth

import (
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
)

func Whoami(g *gin.RouterGroup, Db *sqlx.DB) {
	g.Handle("GET", "/auth/whoami", func(c *gin.Context) {
		// 获取 cookie
		cookie, err := c.Cookie("camp-session")
		var response types.WhoAmIResponse
		if err != nil {
			response.Code = 6
			response.Data = types.TMember{"", "", "", 0}
			c.JSON(http.StatusBadRequest, response)
			return
		}
		var info []types.TMember

		Db.Select(&info, "select * from userinfo where userid=?", cookie)
		response.Code = 0
		response.Data = info[0]
		c.JSON(http.StatusOK, response)
	})
}
