package utils

import (
	"log"
	"os"
	"os/exec"
)

func ToWebP(input string, output string) error {
	command := exec.Command("cwebp", "-q", "85", input, "-o", output)
	command.Stdout = nil
	command.Stderr = nil

	err := command.Run()
	if err != nil {
		log.Printf("%s - error converting: input: %s | output: %s\n", err, input, output)
		return err
	}

	err = os.Remove(input)
	if err != nil {
		return err
	}

	return nil
}
