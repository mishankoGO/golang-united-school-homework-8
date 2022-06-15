package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Users struct {
	Users []User `json:"users"`
}

func (u *Users) Append(usr User) {
	u.Users = append(u.Users, usr)
}

func (u Users) Print(writer io.Writer) {
	var res strings.Builder
	res.WriteString("[")
	for i := 0; i < len(u.Users); i++ {
		res.WriteString("{")
		res.WriteString(fmt.Sprintf(`"id": "%d", "email": "%s", "age": %d`,
			u.Users[i].Id,
			u.Users[i].Email,
			u.Users[i].Age,
		))
		res.WriteString("}")
	}
	res.WriteString("]\n")
	writer.Write([]byte(res.String()))
}

type User struct {
	Id    int
	Email string
	Age   int
}

type Arguments struct {
	Id        int
	Operation string
	Item      User
	FileName  string
}

func stringSplitter(s string) []string {
	//fmt.Println()
	return strings.Split(s, ":")
}

func parseArgs() Arguments {
	op := flag.String("operation", "list", "operation to perform on users list")
	item := flag.String("item", "", "item to add")
	Id := flag.Int("id", -1, "user id")
	fileName := flag.String("fileName", "users.json", "name of the file with users list")
	flag.Parse()

	var args Arguments

	args.Operation = *op
	args.Id = *Id

	if len(*item) != 0 {
		fmt.Println(*item)
		s := strings.Trim(*item, "{}")
		item := strings.Split(s, ",")
		id, err := strconv.Atoi(strings.Trim(stringSplitter(item[0])[1], " "))
		if err != nil {
			fmt.Println("error converting id to int")
		}
		email := strings.Trim(stringSplitter(item[1])[1], " ")
		age, err := strconv.Atoi(strings.Trim(stringSplitter(item[2])[1], " "))
		if err != nil {
			fmt.Println("error converting age to int")
		}
		args.Item = User{id, email, age}
	}

	if len(*fileName) == 0 {
		fmt.Println("please enter the file name")
	}
	args.FileName = *fileName
	fmt.Println(args)

	return args
}

func Add(item User, fileName string) {
	f, err := os.ReadFile(fileName)
	if os.IsNotExist(err) {
		f, err := os.Create(fileName)
		if err != nil {
			fmt.Println("error creating a file")
		}
		defer f.Close()
	}

	users := Users{}
	_ = json.Unmarshal(f, &users)
	fmt.Println(users)

	users.Append(item)
	fmt.Println(users)
	updUsers, _ := json.Marshal(&users)
	os.WriteFile(fileName, updUsers, 0644)
}

func List(fileName string, writer io.Writer) {
	f, err := os.ReadFile(fileName)
	if os.IsNotExist(err) {
		fmt.Println("no such file")
	}
	users := Users{}
	_ = json.Unmarshal(f, &users)
	fmt.Println(users)
	users.Print(writer)

}

func Remove(id int, fileName string) {
	f, err := os.ReadFile(fileName)
	if os.IsNotExist(err) {
		fmt.Println("no such file")
	}
	users := Users{}
	newUsers := Users{}
	_ = json.Unmarshal(f, &users)
	for i := 0; i < len(users.Users); i++ {
		if users.Users[i].Id != id {
			newUsers.Append(users.Users[i])
		}
	}
	updUsers, _ := json.Marshal(&newUsers)
	os.WriteFile(fileName, updUsers, 0644)
}

func FindById(id int, fileName string) {
	f, err := os.ReadFile(fileName)
	if os.IsNotExist(err) {
		fmt.Println("no such file")
	}
	users := Users{}
	_ = json.Unmarshal(f, &users)
	fmt.Println(users)

}

func Perform(args Arguments, writer io.Writer) error {
	op := args.Operation

	switch op {
	case "add":
		item := args.Item
		fileName := args.FileName
		Add(item, fileName)

	case "list":
		fileName := args.FileName
		List(fileName, writer)

	case "remove":
		fileName := args.FileName
		id := args.Id
		Remove(id, fileName)
	}
	return nil
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
