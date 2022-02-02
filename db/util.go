package db

import (
	"bytedance/types"
)

func GetMemberByID(userID int) (types.TMember, types.ErrNo) {
	var ret types.TMember
	if err := NewDB().Get(&ret, "select member_id,member_name,member_nickname,member_type from member where member_id = ? limit 1", userID); err != nil {
		//获取目前的最大自增键
		var maxID int
		if err := NewDB().Get(&maxID, "select max(member_id) from member"); err != nil {
			return ret, types.UnknownError
		}
		//检查是否已删除
		if userID < maxID {
			return ret, types.UserHasDeleted
		}
		//如果不是已删除，则说明用户不存在
		return ret, types.UserNotExisted
	}
	return ret, types.OK
}
