package globals

import (
	"fmt"

	"github.com/seanch0n/halp/cheats"
)

var TheCheat cheats.CheatSearch
var OUTCMD string

func Set(b string) {
	OUTCMD = b
}

func Get() string {
	return OUTCMD
}

func Pr() {
	fmt.Println(OUTCMD)
}
