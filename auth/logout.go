//no redis
//no sql
package auth

import (
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Logout(c *gin.Context) {
	var response types.LogoutResponse

	//重复登出返回成功

	// 删除 cookie
	c.SetCookie("camp-session", "", -1, "/",
		"", false, true)

	response.Code = types.OK
	c.JSON(http.StatusOK, response)
}
