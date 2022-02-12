package course

import (
	"bytedance/db"
	"bytedance/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func UnBindCourseTeacher(c *gin.Context) {
	var request types.UnbindCourseRequest
	var response types.UnbindCourseResponse

	// 先检查参数是否合法
	if err := c.Bind(&request); err != nil {
		response.Code = types.ParamInvalid
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// -----根据 cookie 判断cookie是否过期 : cookie过期则返回 LoginRequired（用户未登录） ------ //
	_, err := c.Cookie("camp-session")
	if err != nil {
		response.Code = types.LoginRequired // cookie 过期，用户未登录
		c.JSON(http.StatusBadRequest, response)
		return
	}

	intCourseID, _ := strconv.Atoi(request.CourseID)
	intTeacherID, _ := strconv.Atoi(request.TeacherID)
	fmt.Println(intCourseID, intTeacherID, request)
	// -----根据 request 判断课程是否存在 : 课程不存在则返回 CourseNotExisted（课程不存在）   ------ //
	if errNo := db.GetCourse_ScheduleByCourseIDAndTeacherID(intCourseID, intCourseID); errNo != types.OK {
		response.Code = errNo
		c.JSON(http.StatusOK, response)
		return
	}

	result, _ := db.NewDB().Exec("delete from course_schedule where course_id=? and member_id=?", intCourseID, intTeacherID)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		response.Code = types.CourseNotBind
	} else {
		response.Code = types.OK
	}

	c.JSON(http.StatusOK, response)
}
