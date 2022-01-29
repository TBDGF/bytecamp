package member

import (
	"bytedance/config"
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
		c.JSON(http.StatusBadRequest, response)
		return
	}

	config.NewDB().Select(&response.Data.MemberList,
		"select * from userinfo where userid between ? and ?", request.Offset, request.Offset+request.Limit)

	response.Code = types.OK
	c.JSON(http.StatusOK, response)
	return
}
