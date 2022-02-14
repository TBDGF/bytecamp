package redis_server

import (
	"bytedance/db"
	"bytedance/types"
	"encoding/json"
	"log"
	"strconv"
)

func GetKeyOfMember(memberID string) string {
	return "camp:member:" + memberID + ":json"
}

func GetMemberByID(memberID string) (types.TMember, types.ErrNo) {
	result, err := NewClient().Get(GetKeyOfMember(memberID)).Result()
	//缓存无结果，查询数据库
	if err != nil {
		member, errNo := db.GetMemberByID(memberID)
		if errNo != types.OK {
			return member, errNo
		}
		jsonString, _ := json.Marshal(member)
		//加入缓存
		NewClient().Set(GetKeyOfMember(memberID), jsonString, 0)
		return member, types.OK
	}
	log.Println("GetMemberByID in redis")
	var member types.TMember
	json.Unmarshal([]byte(result), &member)
	return member, types.OK
}

func DeleteMemberByID(memberID string) {
	NewClient().Del(GetKeyOfMember(memberID))
}

func GetKeyOfCourseAvail(courseID string) string {
	return "camp:course.avail:" + courseID + ":int"
}

func GetCourseAvailByID(courseID string) (int, types.ErrNo) {
	result, err := NewClient().Get(GetKeyOfCourseAvail(courseID)).Result()
	//缓存无结果，查询数据库
	if err != nil {
		avail, errNo := db.GetCourseAvailByID(courseID)
		if errNo != types.OK {
			return avail, errNo
		}
		//加入缓存
		NewClient().Set(GetKeyOfCourseAvail(courseID), avail, 0)
		return avail, errNo
	}
	log.Println("GetCourseAvailByID in redis")
	avail, _ := strconv.Atoi(result)
	return avail, types.OK
}

func GetKeyOfStudentSchedule(studentID string, courseID string) string {
	return "camp:student:" + studentID + ":course:" + courseID + ":int"
}

func GetStudentSchedule(studentID string, courseID string) (bool, types.ErrNo) {
	_, err := NewClient().Get(GetKeyOfStudentSchedule(studentID, courseID)).Result()
	//缓存无结果，查询数据库
	if err != nil {
		var count int
		if err := db.NewDB().Get(&count,
			"select count(*) from student_schedule where student_id = ? AND course_id = ? limit 1",
			studentID, courseID); err != nil {
			return false, types.UnknownError
		}
		//判断是否不存在
		if count == 0 {
			return false, types.StudentHasNoCourse
		}
		//加入缓存
		NewClient().Set(GetKeyOfStudentSchedule(studentID, courseID), 1, 0)
		return true, types.OK
	}
	log.Println("GetStudentSchedule in redis")
	return true, types.OK
}
