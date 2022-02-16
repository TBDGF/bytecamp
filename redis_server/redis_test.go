package redis_server

import (
	"bytedance/db"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	db.InitDB()
	InitRedis()

	log.Println(GetMemberByID("1"))
}

func TestTxDecr(t *testing.T) {
	db.InitDB()
	InitRedis()

	studentAmount := 50000
	courseAmount := 11
	var wg sync.WaitGroup
	wg.Add(studentAmount)

	for i := 1; i <= studentAmount; i++ {
		go func(id int) {
			defer wg.Done()

			studentIDString := strconv.Itoa(id)
			courseIDString := strconv.Itoa(rand.Intn(courseAmount))

			for {
				err := NewClient().Watch(TxDecr(GetKeyOfCourseAvail(courseIDString)),
					GetKeyOfCourseAvail(courseIDString))
				if err == nil {
					log.Println(id, "OK")
					//继续向下执行
					break
				} else if err.Error() == "CourseNotAvailable" {
					log.Println(id, err)
					return
				} else if err.Error() == "CourseNotExisted" {
					log.Println(id, err)
					return
				} else {
					time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)
				}
			}
			////redis添加关系
			//NewClient().Set(GetKeyOfStudentSchedule(studentIDString, courseIDString), 1, 0)

			// --- 更新数据库 --- //
			db.NewDB().Exec("update course set course_available = course_available-1 where course_id = ?", courseIDString)
			db.NewDB().Exec("INSERT INTO camp.student_schedule (student_id, course_id) VALUES (?, ?);", studentIDString, courseIDString)
		}(i)
	}
	wg.Wait()
}
