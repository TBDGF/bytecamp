package course

import (
	"bytedance/db"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Create(c *gin.Context) {
	var request types.CreateCourseRequest
	var response types.CreateCourseResponse

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

	result, _ := db.NewDB().Exec("INSERT INTO camp.course (course_name, course_available) VALUES (?, ?);", request.Name, request.Cap)
	CourseID, _ := result.LastInsertId()
	response.Code = types.OK
	response.Data.CourseID = strconv.Itoa(int(CourseID))
	c.JSON(http.StatusOK, response)

}
