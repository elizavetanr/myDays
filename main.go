package main

import (
	"fmt"
	"github.com/elizavetanr/myDays/calendar"
	"github.com/elizavetanr/myDays/cmd"
	"github.com/elizavetanr/myDays/logger"
	"github.com/elizavetanr/myDays/storage"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	s := storage.NewJsonStorage("calendar.json")
	c := calendar.NewCalendar(s)

	err := c.Load()
	if err != nil {
		fmt.Println("Ошибка: ", err)
	}
	err = logger.Init()
	if err != nil {
		fmt.Println("Ошибка: ", err)
	}
	cli := cmd.NewCmd(c)
	cli.Run()
}
