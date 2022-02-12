package course

import (
	"bytedance/db"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
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

	course, errNo := db.GetCourseByID(request.CourseID)
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

	_, err := db.NewDB().Exec("INSERT INTO camp.teacher_schedule (teacher_id, course_id) VALUES (?, ?);", request.TeacherID, request.CourseID)
	if err != nil {
		errNo = types.UnknownError
	} else {
		errNo = types.OK
	}
	response.Code = errNo

	c.JSON(http.StatusOK, response)
}
