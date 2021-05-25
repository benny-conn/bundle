package term

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
)

func Print(s interface{}) {
	fmt.Printf("%s %s", Blue("[Bundle]"), s)
}
func Println(s interface{}) {
	fmt.Printf("%s %s\n", Blue("[Bundle]"), s)
}
