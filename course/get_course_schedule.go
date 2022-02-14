package course

import (
	"bytedance/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func defn(N int) {
	h = make([]int, N, N)

	match = make([]int, N, N)
	st = make([]bool, N, N)

	key = make([]string, N, N)
	val = make([]string, N, N)
}

func defm(M int) {
	e = make([]int, M, M)
	ne = make([]int, M, M)
}

// 邻接表
var (
	e   []int
	ne  []int
	h   []int
	key []string
	val []string
	idx = 0
)

// 邻接表插入
func insert(a, b int) {
	idx++
	e[idx] = b
	ne[idx] = h[a]
	h[a] = idx
}

var match []int // 匹配状态， 课程匹配了 哪一个老师
var st []bool   // 一轮状态

func find(x int) bool {
	for i := h[x]; i > 0; i = ne[i] {
		t := e[i]
		if !st[t] {
			st[t] = true                         // 这一轮，这个课程 t 被标记匹配了
			if match[t] == 0 || find(match[t]) { // 尝试让原来找到匹配的老师重新匹配其他课程（同一轮中）
				match[t] = x // 课程t匹配了老师x
				return true
			}
		}
	}
	return false
}

func GetCourseSchedule(c *gin.Context) {

	var request types.ScheduleCourseRequest
	var response types.ScheduleCourseResponse
	if err := c.Bind(&request); err != nil {
		response.Code = types.ParamInvalid
		c.JSON(http.StatusOK, response)
		return
	}

	//var teacherAmount, courseAmount, scheduleAmount int // 老师数量、课程数量 关系数量
	teacherAmount := 1
	courseAmount := 1
	scheduleAmount := 1
	mp := make(map[string]int)

	lens := len(request.TeacherCourseRelationShip) + 1
	key = make([]string, lens, lens)
	for k, v := range request.TeacherCourseRelationShip {
		key[teacherAmount] = k
		for i := 0; i < len(v); i++ {
			if mp[v[i]] == 0 {
				mp[v[i]] = courseAmount
				//fmt.Println(courseAmount)
				courseAmount++
			}
			scheduleAmount++
		}
		teacherAmount++
	}

	//初始化全局变量
	if teacherAmount > courseAmount {
		defn(teacherAmount)
	} else {
		defn(courseAmount)
	}
	defm(scheduleAmount)
	idx = 0

	teacherAmount = 1
	courseAmount = 1
	scheduleAmount = 1
	mp = make(map[string]int)
	for k, v := range request.TeacherCourseRelationShip {
		//fmt.Println(k, v)
		key[teacherAmount] = k
		for i := 0; i < len(v); i++ {
			if mp[v[i]] == 0 {
				mp[v[i]] = courseAmount
				val[courseAmount] = v[i]
				courseAmount++
			}
			insert(teacherAmount, mp[v[i]])
			scheduleAmount++
		}
		teacherAmount++
	}

	//fmt.Println(teacherAmount, courseAmount, scheduleAmount)

	teacherAmount--
	courseAmount--
	scheduleAmount--

	res := 0

	// 循环对每一个老师进行匹配
	for i := 1; i <= teacherAmount; i++ {
		for j := 1; j <= courseAmount; j++ { // 重置一轮的状态
			st[j] = false
		}
		if find(i) {
			res++
		}
	}

	response.Data = make(map[string]string)
	for i := 1; i <= courseAmount; i++ {
		if match[i] != 0 {
			response.Data[key[match[i]]] = val[i]
		}
	}
	response.Code = types.OK
	c.JSON(http.StatusOK, response)
}
