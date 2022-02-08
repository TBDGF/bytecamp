package course

import (
	"bytedance/db"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//// 老师绑定课程
//// Method： Post
//// 注：这里的 teacherID 不需要做已落库校验
//// 一个老师可以绑定多个课程 , 不过，一个课程只能绑定在一个老师下面
//type BindCourseRequest struct {
//	CourseID  string `form:"course_id"`
//	TeacherID string `form:"member_id"`
//}
//
//type BindCourseResponse struct {
//	Code ErrNo
//}

func BindCourseTeacher(c *gin.Context) {
	var request types.BindCourseRequest
	var response types.BindCourseResponse

	if err := c.Bind(&request); err != nil {
		response.Code = types.ParamInvalid
		c.JSON(http.StatusBadRequest, response)
		return
	}
	intCourseID, _ := strconv.Atoi(request.CourseID)

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

	course, errNo := db.GetCourseByID(intCourseID)
	if errNo != types.OK {
		response.Code = errNo
		c.JSON(http.StatusBadRequest, response)
		return
	}

	//判断是否重复绑定
	if course.TeacherID != "" {
		response.Code = types.CourseHasBound
		c.JSON(http.StatusBadRequest, response)
		return
	}

	//这里绑定的是老师，所以member_type预定义为3
	_, err = db.NewDB().Exec("INSERT INTO camp.course_schedule (course_id, member_id, member_type) VALUES (?, ?, ?);", request.CourseID, request.TeacherID, 3)
	if err != nil {
		errNo = types.UnknownError
	} else {
		errNo = types.OK
	}
	response.Code = errNo

	c.JSON(http.StatusOK, response)
}
