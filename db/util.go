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
	var result types.TCourse
	if err := NewDB().Get(&result, "select c.course_id,c.course_name,cs.member_id from course c left join course_schedule cs on c.course_id=cs.course_id and cs.member_type=3 where c.course_id = ?;", courseID); err != nil {
		////未绑定
		//if result.CourseID != "" {
		//	return result, types.OK
		//}
		//课程不存在
		return result, types.CourseNotExisted
	}
	return result, types.OK
}

// 判断courseID是否在course_schedule这个表中，如果不在就说明未绑定，如果在就说明已经绑定，因为course_schedule表中，课程对老师的关系是一对一的关系
func GetCourse_ScheduleByID(courseID int) types.ErrNo {
	var result types.TCourse
	if err := NewDB().Get(&result, "select course_id from course_schedule where course_id = ?;", courseID); err != nil {
		// 未绑定
		if result.CourseID == "" {
			return types.OK
		}
	}
	return types.CourseHasBound
}

func GetCourse_ScheduleByCourseIDAndTeacherID(courseID int, teacherID int) types.ErrNo {
	var result types.TCourse
	if err := NewDB().Get(&result, "select course_id,member_id from course_schedule where course_id = ? and member_id=?;", courseID, teacherID); err != nil {
		// 未绑定
		if result.CourseID == "" {
			return types.CourseNotBind
		}
	}
	return types.OK
}
