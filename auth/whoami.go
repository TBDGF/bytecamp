package auth

import (
	"bytedance/db"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

	intID, _ := strconv.Atoi(cookie)
	ret, errNo := db.GetMemberByID(intID)
	response.Code = errNo
	response.Data = ret
	c.JSON(http.StatusOK, response)
}
