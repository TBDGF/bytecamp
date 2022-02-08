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

	//ParamInvalid       ErrNo = 1  // 参数不合法
	//UserHasExisted     ErrNo = 2  // 该 Username 已存在
	//UserHasDeleted     ErrNo = 3  // 用户已删除
	//UserNotExisted     ErrNo = 4  // 用户不存在
	//WrongPassword      ErrNo = 5  // 密码错误
	//LoginRequired      ErrNo = 6  // 用户未登录
	//CourseNotAvailable ErrNo = 7  // 课程已满
	//CourseHasBound     ErrNo = 8  // 课程已绑定过
	//CourseNotBind      ErrNo = 9  // 课程未绑定过
	//PermDenied         ErrNo = 10 // 没有操作权限
	//StudentNotExisted  ErrNo = 11 // 学生不存在
	//CourseNotExisted   ErrNo = 12 // 课程不存在
	//StudentHasNoCourse ErrNo = 13 // 学生没有课程
	//StudentHasCourse   ErrNo = 14 // 学生有课程
	//
	//UnknownError ErrNo = 255 // 未知错误

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
