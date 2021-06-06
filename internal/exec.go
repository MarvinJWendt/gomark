package internal

import (
	"fmt"
	"os/exec"
)

func GetGoDoc(pkgPath string) (GoDoc, error) {
	output, err := exec.Command("go", "doc", "-all", pkgPath).CombinedOutput()
	if err != nil {
		return GoDoc{Raw: string(output)}, fmt.Errorf(`error while running "go doc":` + "\n" + string(output))
	}

	return GoDoc{Raw: string(output)}, nil
}
