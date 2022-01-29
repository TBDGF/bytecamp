package auth

import (
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
)

func Login(g *gin.RouterGroup, Db *sqlx.DB) {

	g.Handle("POST", "/auth/login", func(c *gin.Context) {
		var response types.LoginResponse
		var form types.LoginRequest
		if err := c.Bind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var psd []string
		Db.Select(&psd, "select password from users where name=?", form.Username)
		if len(psd) != 0 && psd[0] == form.Password {
			response.Code = types.OK
			var id []string
			Db.Select(&id, "select userid from userinfo where username=?", form.Username)
			response.Data.UserID = id[0]
			// 设置 cookie
			c.SetCookie("camp-session", id[0], 3000, "/",
				"127.0.0.1", false, true)
			c.JSON(http.StatusOK, response)
		} else {
			response.Code = types.WrongPassword
			response.Data.UserID = ""
			c.JSON(http.StatusBadRequest, response)
		}
	})
}
