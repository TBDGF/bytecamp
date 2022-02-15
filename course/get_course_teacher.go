package course

import (
	"bytedance/db"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetCourseTeacher(c *gin.Context) {
	// 不限用户，都可以根据member_id来查询课程
	var request types.GetTeacherCourseRequest
	var response types.GetTeacherCourseResponse

	if err := c.Bind(&request); err != nil {
		response.Code = types.ParamInvalid
		fail(&response, c, err)
		return
	}
	db.NewDB().Select(&response.Data.CourseList, "select c.course_id, c.course_name, ts.teacher_id from course c left join teacher_schedule ts on c.course_id = ts.course_id where teacher_id = ?;", request.TeacherID)
	response.Code = types.OK
	c.JSON(http.StatusOK, response)
	return
}
