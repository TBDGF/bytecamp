package course

import (
	"bytedance/db"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetCourse(c *gin.Context) {
	// 根据courseID来查询课程
	var request types.GetCourseRequest
	var response types.GetCourseResponse

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
	intID, _ := strconv.Atoi(request.CourseID)
	ret, errNo := db.GetCourseByID(intID)
	response.Code = errNo
	response.Data = ret

	c.JSON(http.StatusOK, response)
}
