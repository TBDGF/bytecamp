package course

import (
	"bytedance/db"
	"bytedance/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Return_paramInvalid(response *types.CreateCourseResponse, c *gin.Context) {
	response.Code = types.ParamInvalid
	response.Data.CourseID = ""
	c.JSON(http.StatusBadRequest, response)
	return
}

func Return_paramInvalid2(response *types.CreateMemberResponse, c *gin.Context) {
	response.Code = types.ParamInvalid
	response.Data.UserID = ""
	c.JSON(http.StatusBadRequest, response)
	return
}

func Create(c *gin.Context) {
	// 暂时设定只有系统管理员才有权限创建课程
	//var CMrequest types.CreateMemberRequest
	var CMresponse types.CreateMemberResponse
	var CCrequest types.CreateCourseRequest
	var CCresponse types.CreateCourseResponse

	if err := c.Bind(&CCrequest); err != nil {
		Return_paramInvalid(&CCresponse, c)
		return
	}
	// -----验证操作权限 : 无权限返回 PermDenied ------ //
	// 根据 cookie 获取当前用户权限
	cookie, err := c.Cookie("camp-session")
	fmt.Printf("%#v\n", cookie)
	if err != nil {
		CCresponse.Code = types.LoginRequired // cookie 过期，用户未登录
		//CCresponse.Data.UserID = ""
		c.JSON(http.StatusBadRequest, CCresponse)
		return
	}
	var usertype types.UserType
	//An error is returned if the result set is empty.
	if err := db.NewDB().Get(&usertype, "select member_type from member where member_id=?", cookie); err != nil {
		//获取目前的最大自增键
		var maxID int
		if err := db.NewDB().Get(&maxID, "select max(member_id) from member"); err != nil {
			CMresponse.Code = types.UnknownError
			CMresponse.Data.UserID = ""
			c.JSON(http.StatusBadRequest, CMresponse)
			return
		}
		//检查是否已删除
		if intID, _ := strconv.Atoi(cookie); intID < maxID {
			CMresponse.Code = types.UserHasDeleted
			CMresponse.Data.UserID = ""
			c.JSON(http.StatusBadRequest, CMresponse)
			return
		}
		//如果不是已删除，则说明用户不存在
		CMresponse.Code = types.UserNotExisted
		CMresponse.Data.UserID = ""
		c.JSON(http.StatusBadRequest, CMresponse)
		return
	}
	if usertype != types.Admin {
		CMresponse.Code = types.PermDenied
		CMresponse.Data.UserID = ""
		c.JSON(http.StatusBadRequest, CMresponse)
		return
	}
	// --- 验证用户名是否存在, 错误返回UserHasExisted --- //
	var count int
	if err := db.NewDB().Get(&count, "select count(*) from course where course_name = ? limit 1", CCrequest.Name); err != nil {
		return
	}
	fmt.Println("return count:", count)
	if count != 0 {
		CCresponse.Code = types.UserHasExisted
		CCresponse.Data.CourseID = ""
		c.JSON(http.StatusBadRequest, CCresponse)
		return
	}
	// 通过了上面的验证，就创建课程
	if err := c.Bind(&CCrequest); err != nil {
		Return_paramInvalid(&CCresponse, c)
		return
	}
	result, _ := db.NewDB().Exec("INSERT INTO camp.course (course_name, course_available) VALUES (?, ?);", CCrequest.Name, CCrequest.Cap)
	CourseID, _ := result.LastInsertId()
	CCresponse.Code = types.OK
	CCresponse.Data.CourseID = strconv.Itoa(int(CourseID))
	c.JSON(http.StatusOK, CCresponse)

}
