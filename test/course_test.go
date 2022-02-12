package test

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
)

func TestCourses(t *testing.T) {
	file, err := os.Create("courses.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	file.WriteString("Name,Cap\n")
	for i := 0; i < 1000; i++ {
		fmt.Fprintln(file, randomString(5), ",", rand.Int()%64)
	}
}
