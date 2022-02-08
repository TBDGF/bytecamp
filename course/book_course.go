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

func success(response *types.BookCourseResponse, c *gin.Context) {
	response.Code = types.OK
	c.JSON(http.StatusOK, response)
}
func fail(response *types.BookCourseResponse, code types.ErrNo, c *gin.Context) {
	response.Code = code
	c.JSON(http.StatusBadRequest, response)
}

func BookCourse(c *gin.Context) {
	var request types.BookCourseRequest
	var response types.BookCourseResponse
	if err := c.Bind(&request); err != nil {
		response.Code = types.ParamInvalid
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// --- 验证是否课程已绑定过, 错误返回CourseHasBound --- //
	var count int
	if err := db.NewDB().Get(&count, "select count(*) from course_schedule where member_type = ? AND member_id = ? AND course_id = ? limit 1", types.Student, request.StudentID, request.CourseID); err != nil {
		response.Code = types.UnknownError
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if count != 0 {
		fail(&response, types.CourseHasBound, c)
		return
	}

	// --- 验证课程容量是否充足, 错误返回CourseNotAvailable --- //

	var available int
	availableString, err := redis_server.Redis().Get(request.CourseID).Result()

	if err != nil { // 缓存中没有数据
		if err := db.NewDB().Get(&available, "select course_available from course where course_id = ? limit 1", request.CourseID); err != nil {
			fail(&response, types.UnknownError, c)
			return
		}
		err = redis_server.Redis().Set(request.CourseID, available, 0).Err()
		if err != nil {
			fail(&response, types.UnknownError, c)
			return
		}
	} else {
		available, err = strconv.Atoi(availableString)
		if err != nil {
			fail(&response, types.UnknownError, c)
			return
		}
	}
	if available <= 0 {
		fail(&response, types.CourseNotAvailable, c)
		return
	}
	available_int64, err := redis_server.Redis().Decr(request.CourseID).Result()
	available = int(available_int64)
	if err != nil {
		fail(&response, types.UnknownError, c)
		return
	}
	log.Println("available:", available)

	db.NewDB().Exec("update course set course_available = ? where course_id = ?", available, request.CourseID)
	db.NewDB().Exec("insert into course_schedule values(null , ?, ?, ?)", request.CourseID, request.StudentID, types.Student)
	response.Code = types.OK
	c.JSON(http.StatusOK, response)
}
