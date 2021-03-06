package member

import (
	"bytedance/db"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetMemberList(c *gin.Context) {
	var request types.GetMemberListRequest
	var response types.GetMemberListResponse

	err := c.Bind(&request)
	if err != nil {
		response.Code = types.ParamInvalid
		failFmt(response, c, err)
		return
	}

	db.NewDB().Select(&response.Data.MemberList,
		"select member_id,member_name,member_nickname,member_type from member where is_deleted=0 and member_id >=? limit ?", request.Offset, request.Limit)

	response.Code = types.OK
	c.JSON(http.StatusOK, response)
	return
}
