package main

import (
	"./cmd"
	"./db"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
)

func main() {
	home, _ := homedir.Dir()
	dbPath := filepath.Join(home, "tasks.db")

	must(db.Init(dbPath))
	must(cmd.RootCmd.Execute())
}

func must(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
