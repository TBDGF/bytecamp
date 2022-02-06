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
