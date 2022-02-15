//redis
//sql not optimized
//使用乐观锁中的版本控制，
package course

import (
	"bytedance/db"
	"bytedance/redis_server"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func success(response *types.BookCourseResponse, c *gin.Context) {
	response.Code = types.OK
	c.JSON(http.StatusOK, response)
}

func BookCourse(c *gin.Context) {
	var request types.BookCourseRequest
	var response types.BookCourseResponse
	if err := c.Bind(&request); err != nil {
		response.Code = types.ParamInvalid
		fail(&response, c, err)
		return
	}

	// --- 验证学生是否存在, 不存在返回StudentNotExisted --- //
	//redis
	ret, errNo := redis_server.GetMemberByID(request.StudentID)
	if errNo == types.UserNotExisted || errNo == types.UserHasDeleted || ret.UserType != types.Student {
		response.Code = types.StudentNotExisted
		fail(&response, c)
		return
	}

	// --- 验证是否课程已绑定过, 错误返回StudentHasCourse --- //
	if result, _ := redis_server.GetStudentSchedule(request.StudentID, request.CourseID); result == true {
		response.Code = types.StudentHasCourse
		fail(&response, c)
		return
	}

	//乐观锁完成检查并自减
	for {
		err := redis_server.NewClient().Watch(redis_server.TxDecr(redis_server.GetKeyOfCourseAvail(request.CourseID)), redis_server.GetKeyOfCourseAvail(request.CourseID))
		if err == nil {
			log.Println(request.CourseID, "OK")
			//继续向下执行
			break
		} else if err.Error() == "CourseNotAvailable" {
			log.Println(request.CourseID, err)
			response.Code = types.CourseNotAvailable
			fail(&response, c)
			return
		} else if err.Error() == "CourseNotExisted" {
			log.Println(request.CourseID, err)
			response.Code = types.CourseNotExisted
			fail(&response, c)
			return
		} else {
			log.Println(err, "retry")
			time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)
		}
	}

	//// --- 验证课程是否存在与课程容量是否充足, 错误返回CourseNotExisted或CourseNotAvailable --- //
	//avail, errNo := redis_server.GetCourseAvailByID(request.CourseID)
	//if errNo != types.OK {
	//	response.Code = errNo
	//	fail(&response, c)
	//	return
	//}
	//
	//if avail <= 0 { // 容量不足
	//	response.Code = types.CourseNotAvailable
	//	fail(&response, c)
	//	return
	//}
	//
	////redis课程容量自减
	//availableInt64, err := redis_server.NewClient().Decr(redis_server.GetKeyOfCourseAvail(request.CourseID)).Result()
	//avail = int(availableInt64)
	//if err != nil {
	//	response.Code = types.UnknownError
	//	fail(&response, c, err)
	//	return
	//}
	success(&response, c)

	//redis添加关系
	redis_server.NewClient().Set(redis_server.GetKeyOfStudentSchedule(request.StudentID, request.CourseID), 1, 0)

	// --- 更新数据库 --- //
	db.NewDB().Exec("update course set course_available = course_available-1 where course_id = ?", request.CourseID)
	db.NewDB().Exec("INSERT INTO camp.student_schedule (student_id, course_id) VALUES (?, ?);", request.StudentID, request.CourseID)
}
