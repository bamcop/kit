package tty

import (
	"fmt"
	"strings"
)

func Clear(input string) {
	strs := strings.Split(input, "\n")
	for i := 0; i < len(strs)-1; i++ {
		fmt.Printf("\033[65535D")
		fmt.Printf("\033[K")
		fmt.Printf("\033[1A")
	}
	fmt.Printf("\033[65535D")
	fmt.Printf("\033[K")
}

func ClearLine() {
	fmt.Printf("\033[65535D")
	fmt.Printf("\033[K")
}
