/*
错误处理
*/

package codepattern

import "fmt"

type Account struct {
	Name   string
	Age    uint8
	Gender uint8

	err error //预留字段，错误信息
}

func (a *Account) CheckName() *Account {
	if a.err != nil {
		return a
	}

	if len(a.Name) < 4 || len(a.Name) > 10 {
		a.err = fmt.Errorf("name length is not right, should between 4 and 10")
	}
	return a
}

func (a *Account) CheckAge() *Account {
	if a.err != nil {
		return a
	}

	if a.Age > 65 {
		a.err = fmt.Errorf("age is over 65")
	}
	return a
}

func (a *Account) CheckGender() *Account {
	if a.err != nil {
		return a
	}

	if a.Gender != 1 && a.Gender != 2 {
		a.err = fmt.Errorf("gender is not right")
	}
	return a
}

func (a *Account) Print() *Account {
	if a.err == nil {
		fmt.Printf("account information, Name=%s, Age=%d, Gender=%d", a.Name, a.Age, a.Gender)
	}

	return a
}

func (a *Account) Error() error {
	return a.err
}
