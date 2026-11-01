package main

import "os"

func main() {
	if err := os.WriteFile("bootstrap-must-not-write", []byte("unexpected"), 0600); err != nil {
		panic(err)
	}
	panic("source mount was writable")
}
