package course

import (
	"bytedance/db"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func BindCourseTeacher(c *gin.Context) {
	var request types.BindCourseRequest
	var response types.BindCourseResponse

	if err := c.Bind(&request); err != nil {
		response.Code = types.ParamInvalid
		c.JSON(http.StatusOK, response)
		return
	}

	course, errNo := db.GetCourseByID(request.CourseID)
	if errNo != types.OK {
		response.Code = errNo
		c.JSON(http.StatusOK, response)
		return
	}

	//判断是否重复绑定
	if course.TeacherID != "" {
		response.Code = types.CourseHasBound
		c.JSON(http.StatusOK, response)
		return
	}

	_, err := db.NewDB().Exec("INSERT INTO camp.teacher_schedule (teacher_id, course_id) VALUES (?, ?);", request.TeacherID, request.CourseID)
	if err != nil {
		errNo = types.UnknownError
	} else {
		errNo = types.OK
	}
	response.Code = errNo

	c.JSON(http.StatusOK, response)
}
