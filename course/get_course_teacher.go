package course

import (
	"bytedance/db"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetCourseTeacher(c *gin.Context) {
	// 可以根据member_id来查询课程
	var request types.GetTeacherCourseRequest
	var response types.GetTeacherCourseResponse

	if err := c.Bind(&request); err != nil {
		response.Code = types.ParamInvalid
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// -----验证cookie是否过期或者用户未登录 : 过期或者未登录就返回 LoginRequired ------ //
	if _, err := c.Cookie("camp-session"); err != nil {
		response.Code = types.LoginRequired // cookie 过期，用户未登录
		c.JSON(http.StatusBadRequest, response)
		return
	}
	intTeacherID, _ := strconv.Atoi(request.TeacherID)
	db.NewDB().Select(&response.Data.CourseList, "select c.course_id,c.course_name,cs.member_id from course c natural join course_schedule cs where c.course_id=cs.course_id and cs.member_id = ? and cs.member_type=3", intTeacherID)
	response.Code = types.OK
	c.JSON(http.StatusOK, response)
	return
}
