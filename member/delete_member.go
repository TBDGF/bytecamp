package member

import (
	"bytedance/db"
	"bytedance/redis_server"
	"bytedance/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Delete(c *gin.Context) {
	var request types.DeleteMemberRequest
	var response types.DeleteMemberResponse

	if err := c.Bind(&request); err != nil {
		fmt.Println(err)
	}

	// 删除成员
	if _, errNo := db.GetMemberByID(request.UserID); errNo != types.OK {
		response.Code = errNo
		c.JSON(http.StatusBadRequest, response)
		return
	}
	redis_server.DeleteMemberByID(request.UserID)
	db.NewDB().Exec("update member set is_deleted=1 where member_id=?", request.UserID)
	response.Code = types.OK
	c.JSON(http.StatusOK, response)
	return
}
