package main

import (
	_ "github.com/go-sql-driver/mysql"
	"go-faker/cmd"
	"os"
)

func main() {
	os.Exit(cmd.Do(os.Stdin, os.Stdout, os.Stderr))
}
