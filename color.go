package main

import "fmt"

func CyanString(s string) string {
	return fmt.Sprintf("\x1b[36m%s\x1b[0m", s)
}

func GreenPrint(s string) {
	fmt.Printf("\x1b[32m%s\x1b[0m", s)
}

func RedPrint(s string) {
	fmt.Printf("\x1b[31m%s\x1b[0m", s)
}

func YellowPrint(s string) {
	fmt.Printf("\x1b[33m%s\x1b[0m", s)
}
