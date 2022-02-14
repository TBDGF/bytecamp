//no redis
//sql optimized
package auth

import (
	"bytedance/db"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	var response types.LoginResponse
	var request types.LoginRequest
	if err := c.Bind(&request); err != nil {
		response.Code = types.WrongPassword
		c.JSON(http.StatusOK, response)
		return
	}
	var ret types.PsdAndId
	if err := db.NewDB().Get(&ret,
		"select member_password, member_id from member where member_name=? and is_deleted = 0 limit 1", request.Username); err != nil {
		response.Code = types.WrongPassword
		c.JSON(http.StatusOK, response)
		return
	}
	if ret.Psd != request.Password {
		response.Code = types.WrongPassword
		c.JSON(http.StatusOK, response)
		return
	}

	response.Code = types.OK
	response.Data.UserID = ret.Id
	// 设置 cookie
	c.SetCookie("camp-session", ret.Id, 3600, "/",
		"", false, true)
	c.JSON(http.StatusOK, response)
}
