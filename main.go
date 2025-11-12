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

var Storage_File_Name = "C:\\TaskedUp\\TaskedUpStorage.json"
var Backup_Storage_File_Name = "C:\\TaskedUp\\TaskedUpStorageBackup.json"

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
func Load_Tasks() ([]Task, error) {
	if !Storage_File_Exists() {
		return []Task{}, nil
	}
	Storage_File_Content, err := os.ReadFile(Storage_File_Name)
	if err != nil || len(Storage_File_Content) == 0 {
		return []Task{}, nil
	}
	var tasks []Task
	if err := json.Unmarshal(Storage_File_Content, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil

}
func Save_Tasks(tasks []Task) error {
	Data_2_Save, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(Storage_File_Name, Data_2_Save, 0644)
}
func Get_Tasks() (string, error) {
	tasks, err := Load_Tasks()
	if err != nil {
		return "", fmt.Errorf("can't read file, %v", err)
	}
	if len(tasks) == 0 {
		return "No registered tasks :)", nil
	}
	var Builder strings.Builder
	for _, task := range tasks {
		fmt.Fprintf(&Builder, "%d. %s - %s\n", task.Index, task.Description, task.Status)
	}
	return Builder.String(), nil
}
func Update_Task(Task_Index int, Task_Status string) (string, error) {
	tasks, err := Load_Tasks()
	if err != nil {
		return "", fmt.Errorf("error loading tasks, %v", err)
	}
	if len(tasks) == 0 {
		return "", fmt.Errorf("No registered tasks")
	}
	Task_Index--
	if Task_Index < 0 || Task_Index >= len(tasks) {
		return "", fmt.Errorf("No such task")
	}
	Old_Status := tasks[Task_Index].Status
	tasks[Task_Index].Status = Task_Status
	if err := Save_Tasks(tasks); err != nil {
		return "", fmt.Errorf("Error saving tasks, %v", err)
	}
	return fmt.Sprintf("Updated task #%d: %s\nOld status - %s\nNew status - %s", Task_Index+1, tasks[Task_Index].Description, Old_Status, Task_Status), nil
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
func Task_Exists(Task_Index int) bool {
	Storage_File_Content, err := os.ReadFile(Storage_File_Name)
	if err != nil || len(Storage_File_Content) == 0 {
		return false
	}
	var tasks []Task
	if err = json.Unmarshal(Storage_File_Content, &tasks); err != nil {
		return false
	}
	if len(tasks) == 0 {
		return false
	}
	Task_Index--
	return Task_Index >= 0 && Task_Index < len(tasks)
}
func Add_Task(Task_Description string, Task_Status string) (string, error) {
	tasks, err := Load_Tasks()
	if err != nil {
		return "", fmt.Errorf("error loading tasks, %v", err)
	}
	NewIndex := 1
	if len(tasks) > 0 {
		NewIndex = tasks[len(tasks)-1].Index + 1
	}

	New_Task := Task{Index: NewIndex, Description: Task_Description, Status: Task_Status}
	tasks = append(tasks, New_Task)
	if err := Save_Tasks(tasks); err != nil {
		return "", fmt.Errorf("error saving tasks, %v", err)
	}
	return fmt.Sprintf("Added task %s\n", New_Task.Description), nil
}
func Remove_Task(Task_Index int) (string, error) {
	tasks, err := Load_Tasks()
	if err != nil {
		return "", fmt.Errorf("can't load tasks, %v", err)
	}
	if len(tasks) == 0 {
		return "", fmt.Errorf("No registered tasks :)")
	}
	Task_Index--
	if Task_Index < 0 || Task_Index >= len(tasks) {
		return "", fmt.Errorf("No such task")
	}
	Task_2_Remove := tasks[Task_Index]
	tasks = append(tasks[:Task_Index], tasks[Task_Index+1:]...)
	for i := range tasks {
		tasks[i].Index = i + 1
	}
	if err := Save_Tasks(tasks); err != nil {
		return "", fmt.Errorf("Can't save tasks, %v", err)
	}
	if len(tasks) == 0 {
		return fmt.Sprintf("Removed task #%d - %s\nCongrats on finishing all of your tasks!", Task_2_Remove.Index, Task_2_Remove.Description), nil
	}
	return fmt.Sprintf("Removed task #%d - %s\n%s", Task_2_Remove.Index, Task_2_Remove.Description, Random_Greeting(Task_2_Remove.Description)), nil
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
	Backup_Storage_File, err := os.Create(Backup_Storage_File_Name)
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
func Convert_To_TXT() {
	Storage_File_Content, err := Get_Tasks()
	if err != nil {
		fmt.Println("Error converting: ", err)
		return
	}
	if err := os.WriteFile("TaskedUpStorageText.txt", []byte(Storage_File_Content), 0644); err != nil {
		fmt.Println("Can't write to file, ", err)
		return
	}
	fmt.Println("Converted to txt :)")
}
func Display_TXT() (string, error) {
	Storage_File_Content, err := os.ReadFile("TaskedUpStorageText.txt")
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("No text file found")
		}
		return "", fmt.Errorf("Can't read file, %v", err)
	}
	return string(Storage_File_Content), nil
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
			Task_List, err := Get_Tasks()
			if err != nil {
				fmt.Println("Error, ", err)
			} else {
				fmt.Println(Task_List)
			}
		case "2":
			Reader := bufio.NewReader(os.Stdin)
			fmt.Print("Ok! Enter the tasks description: ")
			Task_Description, _ := Reader.ReadString('\n')
			Task_Description = strings.TrimSpace(Task_Description)
			fmt.Print("Great, now enter the tasks status: ")
			Task_Status, _ := Reader.ReadString('\n')
			Task_Status = strings.TrimSpace(Task_Status)
			Message, err := Add_Task(Task_Description, Task_Status)
			if err != nil {
				fmt.Println("Error: ", err)
			} else {
				fmt.Println(Message)
			}

		case "3":
			Reader := bufio.NewReader(os.Stdin)
			var Index int
			fmt.Print("OK! Enter task number: ")
			fmt.Scanln(&Index)
			if !Task_Exists(Index) {
				fmt.Println("No such task")
				break
			}
			fmt.Print("Great! Now, enter the tasks new status: ")
			Status, _ := Reader.ReadString('\n')
			Status = strings.TrimSpace(Status)
			Message, err := Update_Task(Index, Status)
			if err != nil {
				fmt.Println("Error, ", err)
			} else {
				fmt.Println(Message)
			}
		case "4":
			var Index int
			fmt.Print("Enter the tasks number: ")
			fmt.Scanln(&Index)
			Message, err := Remove_Task(Index)
			if err != nil {
				fmt.Println("Error, ", err)
			} else {
				fmt.Println(Message)
			}
		case "5":
			var YesOrNo string
			fmt.Print("Are you sure? {y/n}: ")
			fmt.Scanln(&YesOrNo)
			switch YesOrNo {
			case "y":
				Backup_Tasks()
				Remove_All_Tasks()
			case "Y":
				Backup_Tasks()
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
			var YesOrNo string
			fmt.Print("Are you sure? {y/n}: ")
			fmt.Scanln(&YesOrNo)
			switch YesOrNo {
			case "y":
				Restore_Tasks()
			case "Y":
				Restore_Tasks()
			case "n":
				fmt.Println("OK!")
			case "N":
				fmt.Println("OK!")
			default:
				fmt.Println("Invalid choice")
			}
		case "10":
			Running = false
			fmt.Println("Goodbye!")
			time.Sleep(1 * time.Second)
		case "txtc": // Easter egg
			Convert_To_TXT()
		case "txtv":
			Storage_File_Content, err := Display_TXT()
			if err != nil {
				fmt.Println("Error - ", err)
			} else {
				fmt.Println(Storage_File_Content)
			}

		default:
			fmt.Println("Sorry, but ", User_Choice, " is an invalid choice")
		}
	}

}
