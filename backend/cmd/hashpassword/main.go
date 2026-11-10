package main

import (
	"fmt"
	"os"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

func main() {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Fprintln(os.Stderr, "run this command in an interactive terminal")
		os.Exit(1)
	}

	fmt.Print("Enter an initial password (8-12 characters): ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not read password:", err)
		os.Exit(1)
	}
	passwordLength := utf8.RuneCount(password)
	if passwordLength < 8 || passwordLength > 12 {
		fmt.Fprintln(os.Stderr, "password must be between 8 and 12 characters")
		os.Exit(1)
	}

	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not hash password:", err)
		os.Exit(1)
	}
	fmt.Println(string(hash))
}
