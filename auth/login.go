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
	var form types.LoginRequest
	if err := c.Bind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var psd string
	if err := db.NewDB().Get(&psd, "select member_password from member where member_name=?", form.Username); err != nil {
		response.Code = types.UserNotExisted
		response.Data.UserID = ""
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if psd != form.Password {
		response.Code = types.WrongPassword
		response.Data.UserID = ""
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response.Code = types.OK
	var id int
	db.NewDB().Get(&id, "select member_id from member where member_name=?", form.Username)
	response.Data.UserID = strconv.Itoa(id)
	// 设置 cookie
	c.SetCookie("camp-session", strconv.Itoa(id), 3000, "/",
		"127.0.0.1", false, true)
	c.JSON(http.StatusOK, response)
}
