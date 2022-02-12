package test

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"testing"
)

func charAt(s string, index int) string {
	runes := bytes.Runes([]byte(s))
	for i, rune := range runes {
		if i == index {
			return string(rune)
		}
	}
	return ""
}

func randomString(length int) string {
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	maxPos := len(str)
	var ret string
	for i := 0; i < length; i++ {
		ret += charAt(str, int(rand.Float64()*float64(maxPos)))
	}
	return ret
}

func TestMembers(t *testing.T) {
	file, err := os.Create("members.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	file.WriteString("Username,Password,UserType\n")
	for i := 0; i < 10000; i++ {
		fmt.Fprintln(file, randomString(16), ",", randomString(16)+"1", ",", int(math.Round(rand.Float64()))+2)
	}
}
