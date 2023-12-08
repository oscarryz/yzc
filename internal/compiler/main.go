package compiler

import (
	"fmt"
	"os"
	"os/exec"
)

func Run(input string) {

	if err := os.WriteFile("generated/test.go", []byte(input), 777); err != nil {
		fmt.Printf("write error: %q", err)
		return
	}

	cmd := exec.Command("go", "build", "generated/test.go")
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		return
	}

}
