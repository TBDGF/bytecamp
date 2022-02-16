package course

import (
	"bytedance/db"
	"bytedance/redis_server"
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
		failLog(&response, c, err)
		return
	}

	// --- 验证学生是否存在, 不存在返回StudentNotExisted --- //
	//redis
	ret, errNo := redis_server.GetMemberByID(request.StudentID)
	if errNo == types.UserNotExisted || errNo == types.UserHasDeleted || ret.UserType != types.Student {
		response.Code = types.StudentNotExisted
		failLog(&response, c)
		return
	}

	db.NewDB().Select(&response.Data.CourseList,
		"select c.course_id, c.course_name, IFNULL(ts.teacher_id,'') teacher_id from course c left join teacher_schedule ts on c.course_id = ts.course_id left join student_schedule ss on c.course_id = ss.course_id where student_id = ?", request.StudentID)

	//结果为空
	if len(response.Data.CourseList) == 0 {
		response.Code = types.StudentHasNoCourse
		failLog(&response, c)
		return
	}

	c.JSON(http.StatusOK, response)
}
