package course

import (
	"bytedance/db"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

	if err := c.Bind(&request); err != nil {
		response.Code = types.ParamInvalid
		c.JSON(http.StatusBadRequest, response)
		return
	}
	intTeacherID, _ := strconv.Atoi(request.TeacherID)
	db.NewDB().Select(&response.Data.CourseList, "select c.course_id,c.course_name,cs.member_id from course c natural join course_schedule cs where c.course_id=cs.course_id and cs.member_id = ? and cs.member_type=3", intTeacherID)
	response.Code = types.OK
	c.JSON(http.StatusOK, response)
	return
}
