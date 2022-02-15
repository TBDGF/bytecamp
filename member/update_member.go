package member

import (
	"bytedance/db"
	"bytedance/redis_server"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Update(c *gin.Context) {
	var request types.UpdateMemberRequest
	var response types.UpdateMemberResponse

	err := c.Bind(&request)
	if err != nil {
		response.Code = types.ParamInvalid
		fail(response, c, err)
		return
	}

	// ---- 验证用户昵称: 不小于 4 位，不超过 20 位 ----
	if len(request.Nickname) < 4 || len(request.Nickname) > 20 {
		response.Code = types.ParamInvalid
		fail(response, c)
		return
	}
	if err != nil {
		response.Code = types.ParamInvalid
		fail(response, c, err)
		return
	}
	if _, errNo := db.GetMemberByID(request.UserID); errNo != types.OK {
		response.Code = errNo
		fail(response, c)
		return
	}
	redis_server.DeleteMemberByID(request.UserID)
	db.NewDB().Exec("update member set member_nickname = ? where member_id = ?", request.Nickname, request.UserID)
	response.Code = types.OK
	c.JSON(http.StatusOK, response)
	return
}
