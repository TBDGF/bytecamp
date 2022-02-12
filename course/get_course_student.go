package course

import (
	"bytedance/db"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetStudentCourse(c *gin.Context) {
	// 不限用户，都可以根据courseID来查询课程
	var request types.GetStudentCourseRequest
	var response types.GetStudentCourseResponse
	if err := c.Bind(&request); err != nil {
		response.Code = types.ParamInvalid
		c.JSON(http.StatusBadRequest, response)
		return
	}
	db.NewDB().Select(&response.Data.CourseList,
		"select c.course_id, c.course_name, IFNULL(ts.teacher_id,'') teacher_id from course c left join teacher_schedule ts on c.course_id = ts.course_id left join student_schedule ss on c.course_id = ss.course_id where student_id = ?", request.StudentID)

	//结果为空
	if len(response.Data.CourseList) == 0 {
		response.Code = types.StudentHasNoCourse
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusOK, response)
}
