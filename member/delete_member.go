package member

import (
	"bytedance/db"
	"bytedance/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Delete(c *gin.Context) {
	var request types.DeleteMemberRequest
	var response types.DeleteMemberResponse

	if err := c.Bind(&request); err != nil {
		fmt.Println(err)
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
	intCookie, _ := strconv.Atoi(cookie)
	db.NewDB().Get(&usertype, "select member_type from member where member_id = ? limit 1", intCookie)

	if usertype != types.Admin {
		response.Code = types.PermDenied
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// 删除成员
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

	db.NewDB().Exec("delete from member where member_id=?", intID)
	response.Code = types.OK
	c.JSON(http.StatusOK, response)
	return
}
