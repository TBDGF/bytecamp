package course

import (
	"bytedance/db"
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Create(c *gin.Context) {
	// 暂时设定只有系统管理员才有权限创建课程
	var request types.CreateCourseRequest
	var response types.CreateCourseResponse

	if err := c.Bind(&request); err != nil {
		response.Code = types.ParamInvalid
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// -----验证操作权限 : 无权限返回 PermDenied ------ //
	// 根据 cookie 获取当前用户权限
	cookie, err := c.Cookie("camp-session")
	if err != nil {
		response.Code = types.LoginRequired // cookie 过期，用户未登录
		c.JSON(http.StatusBadRequest, response)
		return
	}
	intID, _ := strconv.Atoi(cookie)
	requester, errNo := db.GetMemberByID(intID)
	if errNo != types.OK {
		response.Code = errNo
		c.JSON(http.StatusOK, response)
		return
	}
	if requester.UserType != types.Admin {
		response.Code = types.PermDenied
		c.JSON(http.StatusBadRequest, response)
		return
	}

	result, _ := db.NewDB().Exec("INSERT INTO camp.course (course_name, course_available) VALUES (?, ?);", request.Name, request.Cap)
	CourseID, _ := result.LastInsertId()
	response.Code = types.OK
	response.Data.CourseID = strconv.Itoa(int(CourseID))
	c.JSON(http.StatusOK, response)

}
