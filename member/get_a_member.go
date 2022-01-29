package member

import (
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
)

func GetMember(g *gin.RouterGroup, Db *sqlx.DB) {
	g.Handle("GET", "/member", func(c *gin.Context) {
		var Request types.GetMemberRequest
		if err := c.Bind(&Request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ret []types.TMember
		err := Db.Select(&ret, "select * from userinfo where userid = ?", Request.UserID)
		var response types.GetMemberResponse
		if err != nil || len(ret) == 0 {
			response.Code = types.UserNotExisted
			response.Data = types.TMember{}
			c.JSON(http.StatusBadRequest, response)
			return
		}
		var valid []string
		err = Db.Select(&valid, "select name from users where name = ?", ret[0].Username)
		if err != nil || len(valid) == 0 {
			response.Code = types.UserHasDeleted
			response.Data = types.TMember{}
			c.JSON(http.StatusBadRequest, response)
			return
		}
		response.Code = types.OK
		response.Data = ret[0]
		c.JSON(http.StatusOK, response)
	})
}
