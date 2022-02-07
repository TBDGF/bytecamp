package router

import (
	"bytedance/auth"
	"bytedance/course"
	"bytedance/member"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	g := r.Group("/api/v1")

	// 成员管理
	g.POST("/member/create", member.Create)
	g.GET("/member", member.GetMember)
	g.GET("/member/list", member.GetMemberList)
	g.POST("/member/update", member.Update)
	g.POST("/member/delete", member.Delete)

	// 登录

	g.POST("/auth/login", auth.Login)
	g.POST("/auth/logout", auth.Logout)
	g.GET("/auth/whoami", auth.Whoami)

	// 排课
	g.POST("/course/create", course.Create)
	g.GET("/course/get", course.GetCourse)

	g.POST("/teacher/bind_course", course.BindCourseTeacher)
	g.POST("/teacher/unbind_course", course.UnBindCourseTeacher)
	g.GET("/teacher/get_course", course.GetCourseTeacher)
	g.POST("/course/schedule")

	// 抢课
	g.POST("/student/book_course")
	g.GET("/student/course")

}
