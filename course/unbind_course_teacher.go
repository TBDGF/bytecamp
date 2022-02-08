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
	// 首先验证当前用户是否为教师，只有教师才有权限绑定课程
	var request types.UnbindCourseRequest
	var response types.UnbindCourseResponse

	// 先检查参数是否合法
	if err := c.Bind(&request); err != nil {
		response.Code = types.ParamInvalid
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// -----验证操作权限 : 无权限返回 PermDenied ------ //
	// 根据 cookie 获取当前用户权限
	cookie, err := c.Cookie("camp-session")
	if err != nil {
		response.Code = types.LoginRequired // cookie 过期，用户未登录
		c.JSON(http.StatusBadRequest, response)
		return
	}
	intID, _ := strconv.Atoi(cookie)
	requester, errNo := db.GetMemberByID(intID)
	if errNo != types.OK {
		response.Code = errNo
		c.JSON(http.StatusOK, response)
		return
	}
	if requester.UserType != types.Teacher {
		response.Code = types.PermDenied
		c.JSON(http.StatusBadRequest, response)
		return
	}

	intCourseID, _ := strconv.Atoi(request.CourseID)
	intTeacherID, _ := strconv.Atoi(request.TeacherID)
	fmt.Println(intCourseID, intTeacherID, request)
	result, _ := db.NewDB().Exec("delete from course_schedule where course_id=? and member_id=?", intCourseID, intTeacherID)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		response.Code = types.CourseNotBind
	} else {
		response.Code = types.OK
	}

	c.JSON(http.StatusOK, response)
}
