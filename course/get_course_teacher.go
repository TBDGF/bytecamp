package course

import (
	"bytedance/db"
	"bytedance/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

//// 获取老师下所有课程
//// Method：Get
//type GetTeacherCourseRequest struct {
//	TeacherID string `form:"member_id"`
//}
//
//type GetTeacherCourseResponse struct {
//	Code ErrNo
//	Data struct {
//		CourseList []*TCourse// 指针表示关联性，Tcourse变化，这个值也随之发生变化
//	}
//}

func GetCourseTeacher(c *gin.Context) {
	// 不限用户，都可以根据member_id来查询课程
	var request types.GetTeacherCourseRequest
	var response types.GetTeacherCourseResponse

	err := c.Bind(&request)
	if err != nil {
		response.Code = types.ParamInvalid
		c.JSON(http.StatusBadRequest, response)
		return
	}

	db.NewDB().Select(&response.Data.CourseList,
		"select course_schedule.member_id,course.course_name,course_schedule.course_id from course_schedule,course where course_schedule.course_id=course.course_id and course_schedule.member_id =?", request.TeacherID)
	fmt.Printf("%#v\n", response.Data.CourseList)
	response.Code = types.OK
	c.JSON(http.StatusOK, response)
	return
}
