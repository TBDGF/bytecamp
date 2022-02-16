package member

import (
	"bytedance/db"
	"bytedance/redis_server"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Delete(c *gin.Context) {
	var request types.DeleteMemberRequest
	var response types.DeleteMemberResponse

	if err := c.Bind(&request); err != nil {
		response.Code = types.UnknownError
		failFmt(response, c, err)
		return
	}

	//删除cookie
	// 删除 cookie
	c.SetCookie("camp-session", "", -1, "/",
		"", false, true)

	// 删除成员
	if _, errNo := db.GetMemberByID(request.UserID); errNo != types.OK {
		response.Code = errNo
		failFmt(response, c)
		return
	}
	redis_server.DeleteMemberByID(request.UserID)
	db.NewDB().Exec("update member set is_deleted=1 where member_id=?", request.UserID)
	response.Code = types.OK
	c.JSON(http.StatusOK, response)
	return
}
