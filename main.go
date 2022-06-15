package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
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
		res.WriteString(fmt.Sprintf(`"id": "%s", "email": "%s", "age": %d`,
			u.Users[i].Id,
			u.Users[i].Email,
			u.Users[i].Age,
		))
		res.WriteString("}")
		if i < len(u.Users)-1 {
			res.WriteString(",")

		}
	}
	res.WriteString("]\n")
	writer.Write([]byte(res.String()))
}

var OperationMissing = errors.New("-operation flag has to be specified")
var FileNameMissing = errors.New("-fileName flag has to be specified")
var OperationNotAllowed = func(s string) error {
	return errors.New(fmt.Sprintf("Operation %s not allowed!", s))
}
var ItemMissing = errors.New("-item flag has to be specified")
var IdMissing = errors.New("-id flag has to be specified")
var IdExists = errors.New("Item with id 1 already exists")

type User struct {
	Id    string
	Email string
	Age   int
}

type Arguments struct {
	Id        string
	Operation string
	Item      string
	ItemUsr   User
	FileName  string
}

func stringSplitter(s string) []string {
	return strings.Split(s, ":")
}

func parseArgs() Arguments {
	op := flag.String("operation", "", "operation to perform on users list")
	item := flag.String("item", "", "item to add")
	Id := flag.String("id", "", "user id")
	fileName := flag.String("fileName", "test.json", "name of the file with users list")
	flag.Parse()

	var args Arguments

	args.Operation = *op
	args.Id = *Id
	if len(*item) != 0 {
		s := strings.Trim(*item, "{}")
		item := strings.Split(s, ",")
		id := strings.Trim(stringSplitter(item[0])[1], " ")
		email := strings.Trim(stringSplitter(item[1])[1], " ")
		age, err := strconv.Atoi(strings.Trim(stringSplitter(item[2])[1], " "))
		if err != nil {
			fmt.Println("error converting age to int")
		}
		args.ItemUsr = User{id, email, age}
		args.Item = fmt.Sprintf("[{\"id\" :\"%s\",\"email\" :\"%s\",\"age\" :%d}]", id, email, age)
	}

	if len(*fileName) == 0 {
		fmt.Println("please enter the file name")
	}
	args.FileName = *fileName

	return args
}

func Add(item User, fileName string, writer io.Writer) {
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("error creating file")
		panic(err)
	}
	buf, _ := ioutil.ReadAll(f)
	defer f.Close()

	users := Users{}
	_ = json.Unmarshal(buf, &users)

	for _, elem := range users.Users {
		if elem.Id == item.Id {
			writer.Write([]byte("Item with id 1 already exists"))
		}
	}
	users.Append(item)
	updUsers, _ := json.Marshal(&users)
	os.WriteFile(fileName, updUsers, 0644)
}

func List(fileName string, writer io.Writer) {
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("error creating file")
		panic(err)
	}
	defer f.Close()
	users := Users{}
	buf, _ := ioutil.ReadAll(f)
	_ = json.Unmarshal(buf, &users)
	if len(users.Users) != 0 {
		users.Print(writer)
	}
}

func Remove(id string, fileName string) {
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

func FindById(id string, fileName string, writer io.Writer) {
	f, err := os.ReadFile(fileName)
	if os.IsNotExist(err) {
		fmt.Println("no such file")
	}
	users := Users{}
	_ = json.Unmarshal(f, &users)
	fmt.Println(users)
	for _, elem := range users.Users {
		if elem.Id == id {
			res := fmt.Sprintf("%v\n", elem)
			writer.Write([]byte(res))
		}
	}

}

func Perform(args Arguments, writer io.Writer) (err error) {

	op := args.Operation
	fileName := args.FileName

	fmt.Println(args)

	if op == "" {
		err = fmt.Errorf("%w", OperationMissing)
		return err
	}
	if fileName == "" {
		err = fmt.Errorf("%w", FileNameMissing)
		return err
	}

	switch op {
	case "add":
		item := args.ItemUsr
		//fileName := args.FileName
		fmt.Println(args.Item)
		Add(item, fileName, writer)
		if item.Age == 0 { // think about nil User
			err = fmt.Errorf("%w", ItemMissing)
		}
		return err

	case "list":
		//fileName := args.FileName
		List(fileName, writer)

	case "remove":
		//fileName := args.FileName
		id := args.Id
		if id == "" {
			err = fmt.Errorf("%w", IdMissing)
		}
		Remove(id, fileName)

	case "findById":
		//fileName := args.FileName
		id := args.Id
		FindById(id, fileName, writer)

	default:
		if op != "" {
			err = fmt.Errorf("%w", OperationNotAllowed(op))
		}
	}

	return err
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
