package course

import (
	"bytedance/db"
	"bytedance/redis_server"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func success(response *types.BookCourseResponse, c *gin.Context) {
	response.Code = types.OK
	c.JSON(http.StatusOK, response)
}
func fail(response *types.BookCourseResponse, code types.ErrNo, c *gin.Context) {
	response.Code = code
	c.JSON(http.StatusOK, response)
}

func BookCourse(c *gin.Context) {
	var request types.BookCourseRequest
	var response types.BookCourseResponse
	if err := c.Bind(&request); err != nil {
		fail(&response, types.ParamInvalid, c)
		return
	}

	// --- 验证学生是否存在, 不存在返回StudentNotExisted --- //
	//redis
	ret, errNo := redis_server.GetMemberByID(request.StudentID)
	if errNo == types.UserNotExisted || errNo == types.UserHasDeleted || ret.UserType != types.Student {
		fail(&response, types.StudentNotExisted, c)
		return
	}

	// --- 验证是否课程已绑定过, 错误返回StudentHasCourse --- //
	var count int
	if err := db.NewDB().Get(&count,
		"select count(*) from student_schedule where student_id = ? AND course_id = ? limit 1",
		request.StudentID, request.CourseID); err != nil {
		response.Code = types.UnknownError
		c.JSON(http.StatusOK, response)
		return
	}
	if count != 0 {
		fail(&response, types.StudentHasCourse, c)
		return
	}

	// --- 验证课程是否存在与课程容量是否充足, 错误返回CourseNotExisted或CourseNotAvailable --- //
	avail, errNo := redis_server.GetCourseAvailByID(request.CourseID)
	if errNo != types.OK {
		fail(&response, errNo, c)
		return
	}

	if avail <= 0 { // 容量不足
		fail(&response, types.CourseNotAvailable, c)
		return
	}

	//更新缓存
	availableInt64, err := redis_server.NewClient().Decr(redis_server.GetKeyOfCourseAvail(request.CourseID)).Result()
	avail = int(availableInt64)
	if err != nil {
		fail(&response, types.UnknownError, c)
		return
	}
	success(&response, c)

	// --- 更新数据库 --- //
	db.NewDB().Exec("update course set course_available = ? where course_id = ?", avail, request.CourseID)
	db.NewDB().Exec("INSERT INTO camp.student_schedule (student_id, course_id) VALUES (?, ?);", request.StudentID, request.CourseID)
}
