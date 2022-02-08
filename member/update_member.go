package member

import (
	"bytedance/db"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Update(c *gin.Context) {
	var request types.UpdateMemberRequest
	var response types.UpdateMemberResponse

	err := c.Bind(&request)
	if err != nil {
		response.Code = types.ParamInvalid
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// -----验证操作权限 : 无权限返回 PermDenied ------ //
	// 根据 cookie 获取当前用户权限
	cookie, err := c.Cookie("camp-session")
	if err != nil {
		response.Code = types.LoginRequired // cookie 过期，用户未登录
		c.JSON(http.StatusBadRequest, response)
		return
	}
	var usertype types.UserType
	db.NewDB().Get(&usertype, "select member_type from member where member_id = ? limit 1", cookie)
	if usertype != types.Admin {
		response.Code = types.PermDenied
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// ---- 验证用户昵称: 不小于 4 位，不超过 20 位 ----
	if len(request.Nickname) < 4 || len(request.Nickname) > 20 {
		response.Code = types.ParamInvalid
		c.JSON(http.StatusBadRequest, response)
		return
	}
	intID, err := strconv.Atoi(request.UserID)
	if err != nil {
		response.Code = types.ParamInvalid
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if _, errNo := db.GetMemberByID(intID); errNo != types.OK {
		response.Code = errNo
		c.JSON(http.StatusBadRequest, response)
		return
	}

	db.NewDB().Exec("update member set member_nickname = ? where member_id = ?", request.Nickname, request.UserID)
	response.Code = types.OK
	c.JSON(http.StatusOK, response)
	return
}
