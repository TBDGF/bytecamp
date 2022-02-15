package course

import (
	"bytedance/db"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetCourse(c *gin.Context) {
	// 不限用户，都可以根据courseID来查询课程
	var request types.GetCourseRequest
	var response types.GetCourseResponse

	if err := c.Bind(&request); err != nil {
		response.Code = types.ParamInvalid
		fail(&response, c, err)
		return
	}
	ret, errNo := db.GetCourseByID(request.CourseID)
	response.Code = errNo
	response.Data = ret

	c.JSON(http.StatusOK, response)
}
