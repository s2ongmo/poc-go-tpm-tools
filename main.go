package main

import "fmt"

const Version = "1.0.0"

func main() {
	fmt.Printf("gotpm-poc version %s
", Version)
	fmt.Println("This is a PoC binary simulating google/go-tpm-tools gotpm CLI.")
	fmt.Println("If this binary has been tampered with, you will see evidence below:")
	fmt.Println("INTEGRITY_CHECK=CLEAN")
}
