package auth

import (
	"bytedance/config"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Whoami(c *gin.Context) {
	// 获取 cookie
	cookie, err := c.Cookie("camp-session")
	var response types.WhoAmIResponse
	if err != nil {
		response.Code = types.LoginRequired
		response.Data = types.TMember{"", "", "", 0}
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var info []types.TMember

	config.NewDB().Select(&info, "select * from userinfo where userid=?", cookie)
	response.Code = types.OK
	response.Data = info[0]
	c.JSON(http.StatusOK, response)
}
