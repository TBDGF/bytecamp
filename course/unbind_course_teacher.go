package course

import (
	"bytedance/db"
	"bytedance/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Return_UnBindParamInvalid(response *types.UnbindCourseResponse, c *gin.Context) {
	response.Code = types.ParamInvalid
	c.JSON(http.StatusBadRequest, response)
	return
}

func UnBindCourseTeacher(c *gin.Context) {
	// 首先验证当前用户是否为教师，只有教师才有权限绑定课程
	var CMrequest types.CreateMemberRequest
	var CMresponse types.CreateMemberResponse
	var BCrequest types.UnbindCourseRequest
	var BCresponse types.UnbindCourseResponse

	// 先检查参数是否合法
	if err := c.Bind(&BCrequest); err != nil {
		Return_UnBindParamInvalid(&BCresponse, c)
		return
	}
	// -----验证操作权限 : 无权限返回 PermDenied ------ //
	// 根据 cookie 获取当前用户权限
	cookie, err := c.Cookie("camp-session")
	fmt.Printf("%#v\n", cookie)
	if err != nil {
		BCresponse.Code = types.LoginRequired // cookie 过期，用户未登录
		//CCresponse.Data.UserID = ""
		c.JSON(http.StatusBadRequest, BCresponse)
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
	if usertype != types.Teacher {
		CMresponse.Code = types.PermDenied
		CMresponse.Data.UserID = ""
		c.JSON(http.StatusBadRequest, CMresponse)
		return
	}
	fmt.Printf("%#v\n", usertype)
	// --- 验证用户名是否存在, 错误返回UserHasExisted --- //
	var count int
	if err := db.NewDB().Get(&count, "select count(*) from Member where member_name = ? limit 1", CMrequest.Username); err != nil {
		return
	}
	fmt.Println("return count:", count)
	if count != 0 {
		CMresponse.Code = types.UserHasExisted
		CMresponse.Data.UserID = ""
		c.JSON(http.StatusBadRequest, CMresponse)
		return
	}
	// 通过了上面的验证，才能进行解绑课程。
	var request types.UnbindCourseRequest
	var response types.UnbindCourseResponse
	//var response db.Ccc
	if err := c.Bind(&request); err != nil {
		response.Code = types.ParamInvalid
		c.JSON(http.StatusBadRequest, response)
		return
	}
	intCourseID, _ := strconv.Atoi(request.CourseID)
	intTeacherID, _ := strconv.Atoi(request.TeacherID)
	errNo := db.UnBindCourseTeacherByID(intCourseID, intTeacherID)
	response.Code = errNo

	c.JSON(http.StatusOK, response)
}
