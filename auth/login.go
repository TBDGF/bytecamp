//no redis
//sql optimized
package auth

import (
	"bytedance/db"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Login(c *gin.Context) {
	var response types.LoginResponse
	var request types.LoginRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusOK, response)
		return
	}
	var psd string
	if err := db.NewDB().Get(&psd, "select member_password from member where member_name=? limit 1", request.Username); err != nil {
		response.Code = types.UserNotExisted
		response.Data.UserID = ""
		c.JSON(http.StatusOK, response)
		return
	}
	if psd != request.Password {
		response.Code = types.WrongPassword
		response.Data.UserID = ""
		c.JSON(http.StatusOK, response)
		return
	}

	response.Code = types.OK
	var id int
	db.NewDB().Get(&id, "select member_id from member where member_name=? limit 1", request.Username)
	response.Data.UserID = strconv.Itoa(id)
	// 设置 cookie
	c.SetCookie("camp-session", strconv.Itoa(id), 3600, "/",
		"", false, true)
	c.JSON(http.StatusOK, response)
}
