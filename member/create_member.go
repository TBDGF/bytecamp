package member

import (
	"bytedance/config"
	"bytedance/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Return_paramInvalid(response *types.CreateMemberResponse, c *gin.Context) {
	response.Code = types.ParamInvalid
	response.Data.UserID = ""
	c.JSON(http.StatusBadRequest, response)
	return
}

func Create(c *gin.Context) {
	var request types.CreateMemberRequest
	var response types.CreateMemberResponse
	if err := c.Bind(&request); err != nil {
		Return_paramInvalid(&response, c)
		return
	}
	// -----验证操作权限 : 无权限返回 PermDenied ------ //
	// 根据 cookie 获取当前用户权限
	cookie, err := c.Cookie("camp-session")
	if err != nil {
		response.Code = types.LoginRequired // cookie 过期，用户未登录
		response.Data.UserID = ""
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var usertype []types.UserType
	config.NewDB().Select(&usertype, "select usertype from userinfo where userid = ?", cookie)
	if usertype[0] != types.Admin {
		response.Code = types.PermDenied
		response.Data.UserID = ""
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// ------------- 错误返回参数不合法 -------------- //

	// ---- 验证用户昵称: 不小于 4 位，不超过 20 位 ----
	if len(request.Nickname) < 4 || len(request.Nickname) > 20 {
		Return_paramInvalid(&response, c)
		return
	}

	// ---- 验证用户名 ----
	if len(request.Username) < 8 || len(request.Username) > 20 {
		Return_paramInvalid(&response, c)
		return
	}

	for i := 0; i < len(request.Username); i += 1 {
		char := request.Username[i]
		if char < 'A' || char > 'z' || char > 'Z' && char < 'a' {
			Return_paramInvalid(&response, c)
			return
		}
	}

	// ---- 验证密码 ----
	if len(request.Password) < 8 || len(request.Password) > 20 {
		Return_paramInvalid(&response, c)
		return
	}

	HasBigCase := false
	HasLowCase := false
	HasNum := false

	for i := 0; i < len(request.Password); i += 1 {
		char := request.Password[i]
		if char >= '0' && char <= '9' {
			HasNum = true
		} else if char >= 'a' && char <= 'z' {
			HasLowCase = true
		} else if char >= 'A' && char <= 'Z' {
			HasBigCase = true
		} else {
			Return_paramInvalid(&response, c)
			return
		}
	}

	if !HasNum || !HasLowCase || !HasBigCase {
		Return_paramInvalid(&response, c)
		return
	}

	// ---- 验证用户类型 ----
	if request.UserType != types.Admin && request.UserType != types.Student && request.UserType != types.Teacher {
		Return_paramInvalid(&response, c)
		return
	}

	// --- 验证用户名是否存在, 错误返回UserHasExisted --- //
	var count []int
	config.NewDB().Select(&count, "select count(*) from userinfo where username = ?", request.Username)
	fmt.Println("return count:", count[0])
	if count[0] != 0 {
		response.Code = types.UserHasExisted
		response.Data.UserID = ""
		c.JSON(http.StatusBadRequest, response)
		return
	}

	config.NewDB().Exec("insert into users(name, password) values(?, ?)", request.Username, request.Password)
	config.NewDB().Exec("insert into userinfo(nickname, username, usertype) values(?, ?, ?)", request.Nickname, request.Username, request.UserType)
	var id []string
	config.NewDB().Select(&id, "select userid from userinfo where username = ?", request.Username)
	response.Code = types.OK
	response.Data.UserID = id[0]
	c.JSON(http.StatusOK, response)

}
