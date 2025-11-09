package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Task struct {
	Index       int    `json:"index"`
	Description string `json:"description"`
}

func Show_Start_Menu() {
	fmt.Println("# * # * # * # * # * # * # * # * # * # *")
	fmt.Println("* 1. Show this menu                   #")
	fmt.Println("# 2. Show all registered tasks        *")
	fmt.Println("* 3. Add new task                     #")
	fmt.Println("# 4. Remove existing task             *")
	fmt.Println("* 5. Delete all tasks                 #")
	fmt.Println("# 6. Backup tasks                     *")
	fmt.Println("* 7. Restore tasks                    #")
	fmt.Println("# 8. We'll see                        *")
	fmt.Println("* 10. Exit                            #")
	fmt.Println("# * # * # * # * # * # * # * # * # * # *")
}
func Show_Menu() {
	fmt.Println("1. Show this menu")
	fmt.Println("2. Show all registered tasks")
	fmt.Println("3. Add new task")
	fmt.Println("4. Remove existing task")
	fmt.Println("5. Delete all tasks")
	fmt.Println("6. Backup tasks")
	fmt.Println("7. Restore tasks")
	fmt.Println("8. We'll see")
	fmt.Println("10. Exit")
}
func Create_Storage_File() {
	Storage_File, err := os.Create("TaskedUpStorage.json")
	if err != nil {
		fmt.Println("Can't create storage file, ", err)
		return
	}
	defer Storage_File.Close()
}
func Storage_File_Exists() bool {
	_, err := os.Stat("TaskedUpStorage.json")
	return !errors.Is(err, os.ErrNotExist)
}
func Get_Tasks() string {
	Storage_File_Content, err := os.ReadFile("TaskedUpStorage.json")
	if err != nil {
		return fmt.Sprintf("Can't read storage file, %v", err)
	}
	if len(Storage_File_Content) == 0 {
		return "No registered tasks :)"
	}
	return string(Storage_File_Content)
}
func Get_Last_Task_Index() int {
	Storage_File_Content, err := os.ReadFile("TaskedUpStorage.json")
	if err != nil {
		fmt.Println("Error reading file, ", err)
		return 0
	}
	if len(Storage_File_Content) == 0 {
		return 0
	}
	var tasks []Task
	if err := json.Unmarshal(Storage_File_Content, &tasks); err != nil {
		fmt.Println("Error decoding file, ", err)
		return 0
	}
	if len(tasks) == 0 {
		return 0
	}
	return tasks[len(tasks)-1].Index
}
func Add_Task(Task_Description string) {
	Storage_File_Content, err := os.ReadFile("TaskedUpStorage.json")
	if os.IsNotExist(err) {
		Storage_File_Content = []byte("[]")
		err = nil
	}
	if err != nil {
		fmt.Println("Error reading file, ", err)
		return
	}
	if len(Storage_File_Content) == 0 {
		Storage_File_Content = []byte("[]")
	}
	var tasks []Task
	if err := json.Unmarshal(Storage_File_Content, &tasks); err != nil {
		fmt.Println("Error decoding file, ", err)
		return
	}
	OldIndex := 0
	if len(tasks) > 0 {
		OldIndex = tasks[len(tasks)-1].Index
	}

	New_Task := Task{Index: OldIndex + 1, Description: Task_Description}
	tasks = append(tasks, New_Task)
	Updated_Tasks, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		fmt.Println("Can't encode task, ", err)
		return
	}
	err = os.WriteFile("TaskedUpStorage.json", Updated_Tasks, 0644)
	if err != nil {
		fmt.Println("Can't add task, ", err)
		return
	}
	fmt.Printf("Added task %s\n", New_Task.Description)
}
func Remove_Task(task_index int) {}
func Remove_All_Tasks() {
	Clean := []byte("")
	err := os.WriteFile("TaskedUpStorage.json", Clean, 0644)
	if err != nil {
		fmt.Println(("Can't remove tasks"))
		return
	}
	fmt.Println(("Removed all tasks"))
}
func Backup_Tasks()  {}
func Restore_Tasks() {}

func main() {
	if !Storage_File_Exists() {
		Create_Storage_File()
	}

	Show_Start_Menu()

	var Running bool = true
	var User_Choice string

	for Running {
		fmt.Print("Choose an operation to perform: ")
		fmt.Scanln(&User_Choice)
		switch User_Choice {
		case "1":
			Show_Menu()
		case "2":
			fmt.Println(Get_Tasks())
		case "3":
			Reader := bufio.NewReader(os.Stdin)
			fmt.Print("Ok! Enter the tasks description: ")
			Task_Description, _ := Reader.ReadString('\n')
			Task_Description = strings.TrimSpace(Task_Description)
			Add_Task(Task_Description)
		case "5":
			Remove_All_Tasks()
		case "10":
			Running = false
		default:
			fmt.Println("Sorry, but ", User_Choice, " is an invalid choice")
		}
	}

}
