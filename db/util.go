package db

import (
	"bytedance/types"
	"fmt"
)

type member struct {
	types.TMember
	HasDeleted int `db:"is_deleted"`
}

func GetMemberByID(userID string) (types.TMember, types.ErrNo) {
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

func GetCourseByID(courseID string) (types.TCourse, types.ErrNo) {
	var result types.TCourse
	if err := NewDB().Get(&result, "select c.course_id,c.course_name,ts.teacher_id from course c left join teacher_schedule ts on c.course_id=ts.course_id where c.course_id = ? limit 1", courseID); err != nil {
		//未绑定教师
		if result.CourseID != "" {
			return result, types.OK
		}
		//课程不存在
		return result, types.CourseNotExisted
	}
	return result, types.OK
}

func GetCourseAvailByID(courseID string) (int, types.ErrNo) {
	var result int
	if err := NewDB().Get(&result, "select course_available from course where course_id = ? limit 1;", courseID); err != nil {
		//课程不存在
		return result, types.CourseNotExisted
	}
	return result, types.OK
}
