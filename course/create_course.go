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

	result, _ := db.NewDB().Exec("INSERT INTO camp.course (course_name, course_available) VALUES (?, ?);", request.Name, request.Cap)
	CourseID, _ := result.LastInsertId()
	response.Code = types.OK
	response.Data.CourseID = strconv.Itoa(int(CourseID))
	c.JSON(http.StatusOK, response)

}
