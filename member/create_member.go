package member

import (
	"bytedance/db"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	intID, _ := strconv.Atoi(cookie)
	member, errNo := db.GetMemberByID(intID)
	if errNo != types.OK {
		response.Data.UserID = member.UserID
		response.Code = errNo
		c.JSON(http.StatusOK, response)
		return
	}

	if member.UserType != types.Admin {
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
	var count int
	if err := db.NewDB().Get(&count, "select count(*) from member where member_name = ? limit 1", request.Username); err != nil {
		return
	}
	if count != 0 {
		response.Code = types.UserHasExisted
		response.Data.UserID = ""
		c.JSON(http.StatusBadRequest, response)
		return
	}
	result, _ := db.NewDB().Exec("INSERT INTO camp.member (member_name, member_nickname, member_password, member_type) VALUES (?, ?, ?, ?);", request.Username, request.Nickname, request.Password, request.UserType)
	userID, _ := result.LastInsertId()
	response.Code = types.OK
	response.Data.UserID = strconv.Itoa(int(userID))
	c.JSON(http.StatusOK, response)

}
