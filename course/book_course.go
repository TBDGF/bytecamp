package course

import (
	"bytedance/db"
	"bytedance/redis_server"
	"bytedance/types"
	"github.com/gin-gonic/gin"
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
		fail(&response, types.ParamInvalid, c)
		return
	}

	// --- 验证学生是否存在, 不存在返回StudentNotExisted --- //
	intID, _ := strconv.Atoi(request.StudentID)
	ret, errNo := db.GetMemberByID(intID)
	if errNo == types.UserNotExisted || ret.UserType != types.Student {
		fail(&response, types.StudentNotExisted, c)
		return
	}

	// --- 验证是否课程已绑定过, 错误返回CourseHasBound --- //
	var count int
	if err := db.NewDB().Get(&count,
		"select count(*) from course_schedule where member_type = ? AND member_id = ? AND course_id = ? limit 1",
		types.Student, request.StudentID, request.CourseID); err != nil {
		response.Code = types.UnknownError
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if count != 0 {
		fail(&response, types.CourseHasBound, c)
		return
	}

	// --- 验证课程是否存在与课程容量是否充足, 错误返回CourseNotExisted或CourseNotAvailable --- //

	var available int
	availableString, err := redis_server.Redis().Get(request.CourseID).Result()

	if err != nil { // 缓存中没有数据, 查询数据库并写入缓存
		if err := db.NewDB().Get(&available,
			"select course_available from course where course_id = ? limit 1",
			request.CourseID); err != nil { // 课程不存在
			fail(&response, types.CourseNotExisted, c)
			return
		}
		err = redis_server.Redis().Set(request.CourseID, available, 0).Err()
		if err != nil {
			fail(&response, types.UnknownError, c)
			return
		}
	} else { // 缓存中有数据
		available, err = strconv.Atoi(availableString)
		if err != nil {
			fail(&response, types.UnknownError, c)
			return
		}
	}

	if available <= 0 { // 容量不足
		fail(&response, types.CourseNotAvailable, c)
		return
	}
	availableInt64, err := redis_server.Redis().Decr(request.CourseID).Result() // 更新缓存
	available = int(availableInt64)
	if err != nil {
		fail(&response, types.UnknownError, c)
		return
	}
	// --- 更新数据库 --- //
	db.NewDB().Exec("update course set course_available = ? where course_id = ?", available, request.CourseID)
	db.NewDB().Exec("insert into course_schedule values(null , ?, ?, ?)", request.CourseID, request.StudentID, types.Student)
	success(&response, c)
}
