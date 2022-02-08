package course

import (
	"bytedance/db"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func BindCourseTeacher(c *gin.Context) {
	var request types.BindCourseRequest
	var response types.BindCourseResponse

	if err := c.Bind(&request); err != nil {
		response.Code = types.ParamInvalid
		c.JSON(http.StatusBadRequest, response)
		return
	}
	intCourseID, _ := strconv.Atoi(request.CourseID)

	// -----验证cookie是否过期或者用户未登录 : 过期或者未登录就返回 LoginRequired ------ //
	if _, err := c.Cookie("camp-session"); err != nil {
		response.Code = types.LoginRequired // cookie 过期，用户未登录
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// -----验证post请求的CourseId在course_schedule中是否存在（只检查CourseId,如果有必然已经绑定，因为课程对老师是一对一的关系） : 存在就返回CourseHasBound（课程已绑定）  ------ //
	errNo := db.GetCourse_ScheduleByID(intCourseID)
	if errNo != types.OK {
		response.Code = errNo
		c.JSON(http.StatusBadRequest, response)
		return
	}

	//这里绑定的是老师，所以member_type预定义为3
	if _, err := db.NewDB().Exec("INSERT INTO camp.course_schedule (course_id, member_id, member_type) VALUES (?, ?, ?);", request.CourseID, request.TeacherID, 3); err != nil {
		errNo = types.UnknownError
	} else {
		errNo = types.OK
	}
	response.Code = errNo

	c.JSON(http.StatusOK, response)
}
