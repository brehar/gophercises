package cmd

import (
	"../db"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark a task on your task list as complete",
	Run: func(cmd *cobra.Command, args []string) {
		var ids []int

		for _, arg := range args {
			id, err := strconv.Atoi(arg)

			if err != nil {
				fmt.Println("Failed to parse the argument:", arg)
			} else {
				ids = append(ids, id)
			}
		}

		tasks, err := db.AllTasks()

		if err != nil {
			fmt.Println("Something went wrong:", err.Error())
			os.Exit(1)
		}

		for _, id := range ids {
			if id <= 0 || id > len(tasks) {
				fmt.Println("Invalid task number:", id)
				continue
			}

			task := tasks[id-1]
			err := db.DeleteTask(task.Key)

			if err != nil {
				fmt.Printf("Failed to mark #%d as complete. Error: %s\n", id, err)
			} else {
				fmt.Printf("Marked task #%d as complete.\n", id)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(doCmd)
}
