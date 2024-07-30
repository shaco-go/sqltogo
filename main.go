package main

import (
	"fmt"
	"github.com/shaco-go/sqltogo/internal/modelparse"
)

func main() {
	model, err := modelparse.NewModelParse("./example/courseclass.go", "sys_course_class")
	fmt.Println(model, err)
}
