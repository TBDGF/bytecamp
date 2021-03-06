//no redis
//sql optimized
package auth

import (
	"bytedance/db"
	"bytedance/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func failFmt(response interface{}, c *gin.Context, err ...interface{}) {
	if len(err) > 0 {
		fmt.Println("error:", err[0])
	}
	fmt.Println("error response:", response)
	c.JSON(http.StatusOK, response)
	return
}

func Login(c *gin.Context) {
	var response types.LoginResponse
	var request types.LoginRequest
	if err := c.Bind(&request); err != nil {
		response.Code = types.WrongPassword
		failFmt(&response, c, err)
		return
	}
	var ret types.PsdAndId
	if err := db.NewDB().Get(&ret,
		"select member_password, member_id from member where member_name=? and is_deleted = 0 limit 1", request.Username); err != nil {
		response.Code = types.WrongPassword
		failFmt(&response, c)
		return
	}
	if ret.Psd != request.Password {
		response.Code = types.WrongPassword
		failFmt(&response, c)
		return
	}

	response.Code = types.OK
	response.Data.UserID = ret.Id
	// 设置 cookie
	c.SetCookie("camp-session", ret.Id, 3600, "/",
		"", false, true)
	c.JSON(http.StatusOK, response)
}
