package member

import (
	"bytedance/db"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetMember(c *gin.Context) {
	var request types.GetMemberRequest
	var response types.GetMemberResponse

	if err := c.Bind(&request); err != nil {
		response.Code = types.ParamInvalid
		c.JSON(http.StatusBadRequest, response)
		return
	}
	intID, _ := strconv.Atoi(request.UserID)
	ret, errNo := db.GetMemberByID(intID)
	response.Code = errNo
	response.Data = ret
	c.JSON(http.StatusOK, response)
}
