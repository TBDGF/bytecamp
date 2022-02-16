//redis
//no sql
package auth

import (
	"bytedance/redis_server"
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
		failFmt(&response, c, err)
		return
	}

	ret, errNo := redis_server.GetMemberByID(cookie)
	response.Code = errNo
	response.Data = ret
	c.JSON(http.StatusOK, response)
}
