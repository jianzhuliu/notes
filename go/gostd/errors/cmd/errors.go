package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	if _, err := os.Open("non-existing"); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("file not exist")
		}

		var pathError *os.PathError
		if errors.As(err, &pathError) {
			fmt.Println("as pathError:", pathError.Error())
		}

		wrapError := fmt.Errorf("wrap error %w", err)
		fmt.Println("wrapError:", wrapError)

		unwrapError := errors.Unwrap(wrapError)
		fmt.Println("unwrapError:", unwrapError)
	}
}
