package member

import (
	"bytedance/db"
	"bytedance/redis_server"
	"bytedance/types"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func failFmt(response interface{}, c *gin.Context, err ...interface{}) {
	if len(err) > 0 {
		fmt.Println("error:", err[0])
	}
	fmt.Println("error response:", response)
	c.JSON(http.StatusOK, response)
	return
}

func Create(c *gin.Context) {
	var request types.CreateMemberRequest
	var response types.CreateMemberResponse
	if err := c.Bind(&request); err != nil {
		response.Code = types.ParamInvalid
		failFmt(&response, c, err)
		return
	}
	// -----验证操作权限 : 无权限返回 PermDenied ------ //
	// 根据 cookie 获取当前用户权限
	cookie, err := c.Cookie("camp-session")
	if err != nil {
		response.Code = types.LoginRequired // cookie 过期，用户未登录
		failFmt(&response, c, err)
		return
	}
	requester, errNo := db.GetMemberByID(cookie)
	if errNo != types.OK {
		response.Data.UserID = requester.UserID
		response.Code = errNo
		failFmt(&response, c)
		return
	}

	if requester.UserType != types.Admin {
		response.Code = types.PermDenied
		failFmt(&response, c)
		return
	}

	// ------------- 错误返回参数不合法 -------------- //

	// ---- 验证用户昵称: 不小于 4 位，不超过 20 位 ----
	if len(request.Nickname) < 4 || len(request.Nickname) > 20 {
		response.Code = types.ParamInvalid
		failFmt(&response, c)
		return
	}

	// ---- 验证用户名 ----
	if len(request.Username) < 8 || len(request.Username) > 20 {
		response.Code = types.ParamInvalid
		failFmt(&response, c)
		return
	}

	for i := 0; i < len(request.Username); i += 1 {
		char := request.Username[i]
		if char < 'A' || char > 'z' || char > 'Z' && char < 'a' {
			response.Code = types.ParamInvalid
			failFmt(&response, c)
			return
		}
	}

	// ---- 验证密码 ----
	if len(request.Password) < 8 || len(request.Password) > 20 {
		response.Code = types.ParamInvalid
		failFmt(&response, c)
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
			response.Code = types.ParamInvalid
			failFmt(&response, c)
			return
		}
	}

	if !HasNum || !HasLowCase || !HasBigCase {
		response.Code = types.ParamInvalid
		failFmt(&response, c)
		return
	}

	// ---- 验证用户类型 ----
	if request.UserType != types.Admin && request.UserType != types.Student && request.UserType != types.Teacher {
		response.Code = types.ParamInvalid
		failFmt(&response, c)
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
		failFmt(&response, c)
		return
	}

	//添加到数据库
	result, _ := db.NewDB().Exec("INSERT INTO camp.member (member_name, member_nickname, member_password, member_type) VALUES (?, ?, ?, ?);", request.Username, request.Nickname, request.Password, request.UserType)
	userID, _ := result.LastInsertId()

	response.Code = types.OK
	response.Data.UserID = strconv.Itoa(int(userID))
	//添加到redis
	member := types.TMember{response.Data.UserID, request.Nickname, request.Username, request.UserType}
	jsonString, _ := json.Marshal(member)
	redis_server.NewClient().Set(redis_server.GetKeyOfMember(response.Data.UserID), jsonString, 0)

	c.JSON(http.StatusOK, response)

}
