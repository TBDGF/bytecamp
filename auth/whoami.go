package auth

import (
	"bytedance/db"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Whoami(c *gin.Context) {
	var response types.WhoAmIResponse
	// 获取 cookie
	cookie, err := c.Cookie("camp-session")
	if err != nil {
		response.Code = types.LoginRequired
		c.JSON(http.StatusBadRequest, response)
		return
	}

	ret, errNo := db.GetMemberByID(cookie)
	response.Code = errNo
	response.Data = ret
	c.JSON(http.StatusOK, response)
}
