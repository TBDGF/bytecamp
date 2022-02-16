package course

import (
	"bytedance/db"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UnBindCourseTeacher(c *gin.Context) {
	var request types.UnbindCourseRequest
	var response types.UnbindCourseResponse

	// 先检查参数是否合法
	if err := c.Bind(&request); err != nil {
		response.Code = types.ParamInvalid
		failFmt(&response, c, err)
		return
	}

	result, _ := db.NewDB().Exec("delete from teacher_schedule where teacher_id=? and course_id=? limit 1;", request.TeacherID, request.CourseID)
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		response.Code = types.CourseNotBind
	} else {
		response.Code = types.OK
	}

	c.JSON(http.StatusOK, response)
}
