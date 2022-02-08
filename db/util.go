package db

import (
	"bytedance/types"
	"fmt"
)

type member struct {
	types.TMember
	HasDeleted int `db:"is_deleted"`
}

func GetMemberByID(userID int) (types.TMember, types.ErrNo) {
	var ret member
	if err := NewDB().Get(&ret, "select member_id,member_name,member_nickname,member_type,is_deleted from member where member_id = ? limit 1", userID); err != nil {
		fmt.Println(err.Error())
		//用户不存在
		return ret.TMember, types.UserNotExisted
	}
	//用户已删除
	if ret.HasDeleted == 1 {
		return ret.TMember, types.UserHasDeleted
	}
	return ret.TMember, types.OK
}

func GetCourseByID(courseID int) (types.TCourse, types.ErrNo) {
	//var re course
	//var ret Ccc
	var result types.TCourse
	//var request types.GetCourseRequest
	//var response types.GetCourseResponse
	//NewDB().Get(&re,"select course_id,course_name,course_available from course where course_id = ? limit 1", courseID)
	//fmt.Printf("%#v\n",re)
	if err := NewDB().Get(&result, "select course.course_id,course.course_name,course_schedule.member_id from course,course_schedule where course.course_id=course_schedule.course_id and course.course_id = ? limit 1", courseID); err != nil {
		fmt.Printf("%#v\n", result)
		//获取目前的最大自增键
		var maxID int
		if err := NewDB().Get(&maxID, "select max(course_id) from course"); err != nil {
			return result, types.UnknownError
		}
		//检查是否已删除
		if courseID < maxID {
			return result, types.UserHasDeleted
		}
		//如果不是已删除，则说明用户不存在
		return result, types.UserNotExisted
	}
	return result, types.OK // 需要联合查询
}

//type BindCourseRequest struct {
//	CourseID  string
//	TeacherID string
//}
//
//type BindCourseResponse struct {
//	Code ErrNo
//}
type resultCourse struct {
	Course_id        int
	Course_name      string
	Course_available string
}

type resultCourseSchedule struct {
	Course_id int
}

func BindCourseTeacherByID(courseID int, teacherID int) types.ErrNo {
	// 首先验证当前用户是否为教师，教师才能绑定课程

	// 根据request传来的courseid判断在course表中是否存在，因为绑定课程和老师之前想要先创建好课程
	// 根据request传来的memberid判断member表中是否存在此memberid,因为绑定课程和老师之前，老师的信息一定要在用户表中
	var requestCourse resultCourse
	var requestMember types.TMember
	var response types.BindCourseResponse
	var requestCourseSchedule resultCourseSchedule
	//var ret Ccc
	//NewDB().Get(&resultMember,"select course_id from course where course_id = ? limit 1", courseID)
	//fmt.Printf("%#v\n",resultMember)
	// 判断courseid是否存在于course表中
	if err := NewDB().Get(&requestCourse, "select course_id from course where course_id = ? limit 1", courseID); err != nil {
		//获取目前的最大自增键
		//fmt.Printf("%#v\n", requestCourse)
		var maxID int
		if err := NewDB().Get(&maxID, "select max(course_id) from course"); err != nil {
			return types.UnknownError
		}
		//检查是否已删除
		if courseID < maxID {
			return types.CourseNotExisted
		}
		//如果不是已删除，则说明课程不存在
		return types.CourseNotExisted
	}
	// 判断teacherid是否存在于member表中 ,绑定的权限判断是在bind_course_teacher中判断的
	if err := NewDB().Get(&requestMember, "select member_id, member_type from member where member_type=3 and member_id = ? limit 1", teacherID); err != nil {
		//获取目前的最大自增键
		var maxID int
		if err := NewDB().Get(&maxID, "select max(member_id) from member"); err != nil {
			return types.UnknownError
		}
		//检查是否已删除
		if courseID < maxID {
			return types.UserHasDeleted
		}
		//如果不是已删除，则说明用户不存在
		return types.UserNotExisted
	}
	// 我们这里做了限制。一个课程只能绑定一个老师，一个老师可以有多个课程。所以要检查course_schedule表中的course_id不能重复
	if err := NewDB().Get(&requestCourseSchedule, "select course_id from course_schedule where course_id = ? limit 1", courseID); err != nil {
		// 经过了上面的验证，下面就开始往course_schedule表中插入数据 下面的需要联合查询(或者分步查询)
		NewDB().Exec("INSERT INTO camp.course_schedule (course_id, member_id,member_type) VALUES (?, ?, ?);", courseID, teacherID, requestMember.UserType)
		response.Code = types.OK
		return types.OK
	}
	response.Code = types.CourseHasBound
	return types.CourseHasBound
}

//// Method： Post
//type UnbindCourseRequest struct {
//	CourseID  string
//	TeacherID string
//}
//
//type UnbindCourseResponse struct {
//	Code ErrNo
//}
func UnBindCourseTeacherByID(courseID int, teacherID int) types.ErrNo {
	// 首先验证当前用户是否为教师，教师才能解绑课程,解绑的权限判断是在unbind_course_teacher中判断的
	//var request types.UnbindCourseRequest
	var response types.UnbindCourseResponse
	var requestCourseSchedule types.UnbindCourseRequest
	//NewDB().Get(&resultMember,"select course_id from course where course_id = ? limit 1", courseID)
	//fmt.Printf("%#v\n",resultMember)
	// 我们这里做了限制。一个课程只能绑定一个老师，一个老师可以有多个课程。所以要检查course_schedule表中的course_id不能重复
	if err := NewDB().Get(&requestCourseSchedule, "select course_id from course_schedule where course_id = ? limit 1", courseID); err != nil {
		// 经过了上面的验证，下面就开始往course_schedule表中插入数据 下面的需要联合查询(或者分步查询)
		NewDB().Exec("delete from course_schedule where course_id=? and member_id=?", courseID, teacherID)
		response.Code = types.OK
		return types.OK
	}
	response.Code = types.CourseHasBound
	return types.CourseHasBound
}
