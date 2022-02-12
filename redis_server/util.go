package redis_server

import (
	"bytedance/db"
	"bytedance/types"
	"encoding/json"
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
		return member, errNo
	}
	var member types.TMember
	json.Unmarshal([]byte(result), &member)
	return member, types.OK
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
	avail, _ := strconv.Atoi(result)
	return avail, types.OK
}
