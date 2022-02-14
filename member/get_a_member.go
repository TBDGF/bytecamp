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
		c.JSON(http.StatusOK, response)
		return
	}
	ret, errNo := redis_server.GetMemberByID(request.UserID)
	response.Code = errNo
	response.Data = ret
	c.JSON(http.StatusOK, response)
}
