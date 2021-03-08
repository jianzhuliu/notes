package codepattern

import (
	"errors"
	"testing"
)

func checkAccount(a Account) error {
	return a.CheckName().CheckAge().CheckGender().Print().Error()
}

func TestError(t *testing.T) {
	a := Account{Name: "tak", Age: 70, Gender: 0}
	err := checkAccount(a)

	if err == nil {
		t.Fatal("should has name length err")
	}

	if !errors.Is(err, ErrNameLengh) {
		t.Fatal(err)
	}

	a.Name = "take"

	err = checkAccount(a)

	if err == nil {
		t.Fatal("should has age err")
	}

	if !errors.Is(err, ErrAge) {
		t.Fatal(err)
	}

	a.Age = 60

	err = checkAccount(a)

	if err == nil {
		t.Fatal("should has gender err")
	}

	if !errors.Is(err, ErrGender) {
		t.Fatal(err)
	}

	a.Gender = 1

	err = checkAccount(a)

	if err != nil {
		t.Fatal("should has no err")
	}
}
