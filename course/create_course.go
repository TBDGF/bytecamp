//redis
//sql optimized
package course

import (
	"bytedance/db"
	"bytedance/redis_server"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func fail(response interface{}, c *gin.Context, err ...interface{}) {
	if len(err) > 0 {
		log.Println("error:", err[0])
	}
	log.Println("error response:", response)
	c.JSON(http.StatusOK, response)
	return
}

func Create(c *gin.Context) {
	var request types.CreateCourseRequest
	var response types.CreateCourseResponse

	if err := c.Bind(&request); err != nil {
		response.Code = types.ParamInvalid
		fail(&response, c, err)
		return
	}

	result, _ := db.NewDB().Exec("INSERT INTO camp.course (course_name, course_available) VALUES (?, ?);", request.Name, request.Cap)
	courseID, _ := result.LastInsertId()

	response.Code = types.OK
	response.Data.CourseID = strconv.Itoa(int(courseID))

	//添加到redis
	redis_server.NewClient().Set(redis_server.GetKeyOfCourseAvail(response.Data.CourseID), request.Cap, 0)

	c.JSON(http.StatusOK, response)

}
