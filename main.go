package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Task struct {
	Index       int    `json:"index"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

var Storage_File_Name = "TaskedUpStorage.json"
var Backup_Storage_File_Name = "TaskedUpStorageBackup.json"

func Show_Start_Menu() {
	fmt.Println("# * # * # * # * # * # * # * # * # * # *")
	fmt.Println("* 0. Show operation menu              #")
	fmt.Println("# 1. Show all registered tasks        *")
	fmt.Println("* 2. Add new task                     #")
	fmt.Println("# 3. Update existing task             *")
	fmt.Println("* 4. Remove existing task             #")
	fmt.Println("# 5. Delete all tasks                 *")
	fmt.Println("* 6. Backup tasks                     #")
	fmt.Println("# 7. Restore tasks                    *")
	fmt.Println("* 10. Exit                            #")
	fmt.Println("# * # * # * # * # * # * # * # * # * # *")
}
func Show_Menu() {
	fmt.Println("0. Show operation menu")
	fmt.Println("1. Show all registered tasks")
	fmt.Println("2. Add new task")
	fmt.Println("3. Update existing task")
	fmt.Println("4. Remove existing task")
	fmt.Println("5. Delete all tasks")
	fmt.Println("6. Backup tasks")
	fmt.Println("7. Restore tasks")
	fmt.Println("10. Exit")
}
func Create_Storage_File() {
	Storage_File, err := os.Create(Storage_File_Name)
	if err != nil {
		fmt.Println("Can't create storage file, ", err)
		return
	}
	defer Storage_File.Close()
}
func Create_Backup_Storage_File() {
	Backup_Storage_File, err := os.Create(Backup_Storage_File_Name)
	if err != nil {
		fmt.Println("Can't create backup storage file, ", err)
		return
	}
	defer Backup_Storage_File.Close()
}
func Storage_File_Exists() bool {
	_, err := os.Stat(Storage_File_Name)
	return !errors.Is(err, os.ErrNotExist)
}
func Backup_Storage_File_Exists() bool {
	_, err := os.Stat(Backup_Storage_File_Name)
	return !errors.Is(err, os.ErrNotExist)
}
func Get_Tasks() string {
	Storage_File_Content, err := os.ReadFile(Storage_File_Name)
	if err != nil {
		return fmt.Sprintf("Can't read storage file, %v", err)
	}
	if len(Storage_File_Content) == 0 {
		return "No registered tasks :)"
	}
	var tasks []Task
	err = json.Unmarshal(Storage_File_Content, &tasks)
	if err != nil {
		return fmt.Sprintf("Error getting tasks, %v", err)
	}
	if len(tasks) == 0 {
		return "No registered tasks :)"
	}
	result := ""
	for _, task := range tasks {
		result += fmt.Sprintf("%d. %s - %s\n", task.Index, task.Description, task.Status)
	}
	return result
}
func Update_Task(Task_Index int, Task_Status string) {
	Storage_File_Content, err := os.ReadFile(Storage_File_Name)
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
	Task_Index--
	if Task_Index < 0 || Task_Index >= len(tasks) {
		fmt.Println("No such task")
		return
	}
	Old_Status := tasks[Task_Index].Status
	tasks[Task_Index].Status = Task_Status
	Updated_Tasks, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		fmt.Println("Can't encode task, ", err)
		return
	}
	err = os.WriteFile(Storage_File_Name, Updated_Tasks, 0644)
	if err != nil {
		fmt.Println("Can't update task, ", err)
		return
	}
	fmt.Printf("Updated task #%d: %s\nOld status: %s,\nNew status: %s\n", Task_Index, tasks[Task_Index].Description, Old_Status, Task_Status)
}
func Get_Last_Task_Index() int {
	Storage_File_Content, err := os.ReadFile(Storage_File_Name)
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
func Add_Task(Task_Description string, Task_Status string) {
	Storage_File_Content, err := os.ReadFile(Storage_File_Name)
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

	New_Task := Task{Index: OldIndex + 1, Description: Task_Description, Status: Task_Status}
	tasks = append(tasks, New_Task)
	Updated_Tasks, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		fmt.Println("Can't encode task, ", err)
		return
	}
	err = os.WriteFile(Storage_File_Name, Updated_Tasks, 0644)
	if err != nil {
		fmt.Println("Can't add task, ", err)
		return
	}
	fmt.Printf("Added task %s\n", New_Task.Description)
}
func Remove_Task(Task_Index int) {
	Storage_File_Content, err := os.ReadFile(Storage_File_Name)
	if err != nil {
		fmt.Println("Error reading file, ", err)
		return
	}
	if len(Storage_File_Content) == 0 {
		fmt.Println("No registerd tasks :)")
		return
	}
	var tasks []Task
	if err := json.Unmarshal(Storage_File_Content, &tasks); err != nil {
		fmt.Println("Error decoding file, ", err)
		return
	}
	if len(tasks) == 0 {
		fmt.Println("No registerd tasks :)")
		return
	}
	Task_Index -= 1
	if Task_Index < 0 || Task_Index >= len(tasks) {
		fmt.Println("No such task")
		return
	}
	TaskToRemove := tasks[Task_Index]
	tasks = append(tasks[:Task_Index], tasks[Task_Index+1:]...)
	for i := range tasks {
		tasks[i].Index = i + 1
	}
	Updated_Tasks, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		fmt.Println("Error encoding tasks, ", err)
		return
	}
	if err := os.WriteFile(Storage_File_Name, Updated_Tasks, 0644); err != nil {
		fmt.Println("Error removing task, ", err)
		return
	}
	fmt.Printf("Removed task %d - %s\n", TaskToRemove.Index, TaskToRemove.Description)
}
func Get_Task_By_Index(Task_Index int) string {
	Storage_File_Content, err := os.ReadFile(Storage_File_Name)
	if err != nil {
		return fmt.Sprintf("Error reading file, %v", err)
	}
	if len(Storage_File_Content) == 0 {
		return "No registerd tasks :)"
	}
	var tasks []Task
	if err := json.Unmarshal(Storage_File_Content, &tasks); err != nil {
		return fmt.Sprintf("Error decoding file, %v", err)
	}
	if len(tasks) == 0 {
		return "No registerd tasks :)"
	}
	Task_Index -= 1
	if Task_Index < 0 || Task_Index >= len(tasks) {
		return "No such task"
	}
	return tasks[Task_Index].Description
}
func Random_Greeting(Task_Description string) string {
	Rand_Source := rand.NewSource(time.Now().UnixNano())
	Randomed := rand.New(Rand_Source)
	Random_Number := Randomed.Intn(5)
	switch Random_Number {
	case 0:
		return "Well done!"
	case 1:
		return fmt.Sprintf("Congrats on completing task: %s", Task_Description)
	case 2:
		return "Congratulations :)"
	case 3:
		return "Good job"
	case 4:
		return fmt.Sprintf("Finished %s!", Task_Description)
	default:
		return "Hell yeah!"
	}
}
func Remove_All_Tasks() {
	Clean := []byte("")
	err := os.WriteFile(Storage_File_Name, Clean, 0644)
	if err != nil {
		fmt.Println(("Can't remove tasks"))
		return
	}
	fmt.Println(("Removed all tasks"))
}
func Backup_Tasks() {
	Backup_Storage_File, err := os.Create("TaskedUpStorageBackup.json")
	if err != nil {
		fmt.Println("Error creating backup, ", err)
	}
	if !Storage_File_Exists() {
		fmt.Println("No storage file")
		Create_Storage_File()
		return
	}
	Storage_File_Content, err := os.ReadFile(Storage_File_Name)
	if err != nil {
		fmt.Println("Error reading tasks, ", err)
		return
	}
	if len(Storage_File_Content) == 0 {
		fmt.Println("No registered tasks to back up")
		return
	}
	os.WriteFile(Backup_Storage_File_Name, Storage_File_Content, 0644)
	defer Backup_Storage_File.Close()
	fmt.Println("Backed up tasks :)")
}
func Restore_Tasks() {
	if !Backup_Storage_File_Exists() {
		Create_Backup_Storage_File()
	}
	if !Storage_File_Exists() {
		Create_Storage_File()
	}
	Backup_Storage_File_Content, err := os.ReadFile(Backup_Storage_File_Name)
	if err != nil {
		fmt.Println("Error reading tasks, ", err)
		return
	}
	if len(Backup_Storage_File_Content) == 0 {
		fmt.Println("No registered tasks to restore")
		return
	}
	os.WriteFile(Storage_File_Name, Backup_Storage_File_Content, 0644)
	fmt.Println("Restored tasks :)")

}

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
		case "0":
			Show_Menu()
		case "1":
			fmt.Println(Get_Tasks())
		case "2":
			Reader := bufio.NewReader(os.Stdin)
			fmt.Print("Ok! Enter the tasks description: ")
			Task_Description, _ := Reader.ReadString('\n')
			Task_Description = strings.TrimSpace(Task_Description)
			fmt.Print("Great, now enter the tasks status: ")
			Task_Status, _ := Reader.ReadString('\n')
			Task_Status = strings.TrimSpace(Task_Status)
			Add_Task(Task_Description, Task_Status)
		case "3":
			Reader := bufio.NewReader(os.Stdin)
			var Index int
			fmt.Print("OK! Enter task number: ")
			fmt.Scanln(&Index)
			fmt.Print("Great! Now, enter the tasks new status: ")
			Status, _ := Reader.ReadString('\n')
			Status = strings.TrimSpace(Status)
			Update_Task(Index, Status)
		case "4":
			var Index int
			fmt.Print("Enter the tasks number: ")
			fmt.Scanln(&Index)
			Remove_Task(Index)
			Removed_Description := Get_Task_By_Index(Index)
			fmt.Println(Random_Greeting(Removed_Description))
		case "5":
			var YesOrNo string
			fmt.Print("Are you sure? {y/n}: ")
			fmt.Scanln(&YesOrNo)
			switch YesOrNo {
			case "y":
				Remove_All_Tasks()
			case "Y":
				Remove_All_Tasks()
			case "n":
				fmt.Println("OK!")
			case "N":
				fmt.Println("OK!")
			default:
				fmt.Println("Invalid choice")
			}
		case "6":
			Backup_Tasks()
		case "7":
			Restore_Tasks()
		case "10":
			Running = false
		default:
			fmt.Println("Sorry, but ", User_Choice, " is an invalid choice")
		}
	}

}
