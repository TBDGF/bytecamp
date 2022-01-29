package member

import (
	"bytedance/config"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetMember(c *gin.Context) {
	var Request types.GetMemberRequest
	var response types.GetMemberResponse

	if err := c.Bind(&Request); err != nil {
		response.Code = types.ParamInvalid
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var ret []types.TMember
	err := config.NewDB().Select(&ret, "select * from userinfo where userid = ?", Request.UserID)

	if err != nil || len(ret) == 0 {
		response.Code = types.UserNotExisted
		response.Data = types.TMember{}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var valid []string
	err = config.NewDB().Select(&valid, "select name from users where name = ?", ret[0].Username)
	if err != nil || len(valid) == 0 {
		response.Code = types.UserHasDeleted
		response.Data = types.TMember{}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response.Code = types.OK
	response.Data = ret[0]
	c.JSON(http.StatusOK, response)
}
