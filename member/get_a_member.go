package member

import (
	"bytedance/redis_server"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetMember(c *gin.Context) {
	var request types.GetMemberRequest
	var response types.GetMemberResponse

	if err := c.Bind(&request); err != nil {
		response.Code = types.ParamInvalid
		failFmt(response, c, err)
		return
	}
	ret, errNo := redis_server.GetMemberByID(request.UserID)
	response.Code = errNo
	response.Data = ret
	if errNo != types.OK {
		failFmt(response, c)
	}
	c.JSON(http.StatusOK, response)
}
