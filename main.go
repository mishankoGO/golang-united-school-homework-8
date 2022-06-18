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

var OperationMissing = errors.New("-operation flag has to be specified")
var FileNameMissing = errors.New("-fileName flag has to be specified")
var OperationNotAllowed = func(s string) error {
	return errors.New(fmt.Sprintf("Operation %s not allowed!", s))
}
var ItemMissing = errors.New("-item flag has to be specified")
var IdMissing = errors.New("-id flag has to be specified")
var IdExists = func(s string) string {
	return fmt.Sprintf("Item with id %s already exists", s)
}
var IdNotFound = func(s string) string {
	return fmt.Sprintf("Item with id %s not found", s)
}

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type Arguments struct {
	id        string
	operation string
	item      string
	itemUsr   User
	fileName  string
}

func (a Arguments) ItemParser() User {
	item := a.item
	u := User{}

	item = strings.Trim(item, "[{}]")
	item = strings.Replace(item, " ", "", -1)
	itemLst := strings.Split(item, ",")

	for _, elem := range itemLst {
		kVal := strings.Split(elem, ":")
		kVal[0] = strings.Trim(kVal[0], "\"")
		kVal[1] = strings.Trim(kVal[1], "\"")
		if kVal[0] == "id" {
			u.Id = kVal[1]
		} else if kVal[0] == "email" {
			u.Email = kVal[1]
		} else if kVal[0] == "age" {
			age, err := strconv.Atoi(kVal[1])
			if err != nil {
				fmt.Println("error converting age")
			}
			u.Age = age
		}
	}
	return u
}

func stringSplitter(s string) []string {
	return strings.Split(s, ":")
}

func parseArgs() Arguments {
	op := flag.String("operation", "", "operation to perform on users list")
	item := flag.String("item", "", "item to add")
	Id := flag.String("id", "", "user id")
	fileName := flag.String("fileName", "", "name of the file with users list")
	flag.Parse()

	var args Arguments

	args.operation = *op
	args.id = *Id
	if len(*item) != 0 {
		s := strings.Trim(*item, "{}")
		item := strings.Split(s, ",")
		id := strings.Trim(stringSplitter(item[0])[1], " ")
		email := strings.Trim(stringSplitter(item[1])[1], " ")
		age, err := strconv.Atoi(strings.Trim(stringSplitter(item[2])[1], " "))
		if err != nil {
			fmt.Println("error converting age to int")
		}
		args.item = fmt.Sprintf("[{\"id\" :\"%s\",\"email\" :\"%s\",\"age\" :%d}]", id, email, age)
		//args.itemUsr = args.ItemParser()

	}

	args.fileName = *fileName

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

	var users []User
	_ = json.Unmarshal(buf, &users)
	for _, elem := range users {
		if elem.Id == item.Id {
			writer.Write([]byte(IdExists(elem.Id)))
		}
	}
	if item.Age != 0 {
		users = append(users, item)

	}
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
	buf, _ := ioutil.ReadAll(f)
	if len(buf) != 0 {
		writer.Write(buf)
	}
}

func Remove(id string, fileName string, writer io.Writer) {
	f, _ := os.ReadFile(fileName)

	var users, newUsers []User
	_ = json.Unmarshal(f, &users)
	for i := 0; i < len(users); i++ {
		if users[i].Id != id {
			newUsers = append(newUsers, users[i])
		}
	}

	updUsers, _ := json.Marshal(&newUsers)

	if len(newUsers) == len(users) {
		writer.Write([]byte(IdNotFound(id)))
	}
	os.WriteFile(fileName, updUsers, 0644)
}

func FindById(id string, fileName string, writer io.Writer) {
	f, _ := os.ReadFile(fileName)
	var users []User
	_ = json.Unmarshal(f, &users)
	for _, elem := range users {
		if elem.Id == id {
			res := fmt.Sprintf("{\"id\":\"%s\",\"email\":\"%s\",\"age\":%d}", elem.Id, elem.Email, elem.Age)
			writer.Write([]byte(res))
			return
		}
	}
	writer.Write([]byte(""))
}

func Perform(args Arguments, writer io.Writer) (err error) {

	op := args.operation
	fileName := args.fileName

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
		if len(args.item) != 0 {
			args.itemUsr = args.ItemParser()
		}
		item := args.itemUsr
		Add(item, fileName, writer)
		if item.Age == 0 { // think about nil User
			err = fmt.Errorf("%w", ItemMissing)
		}
		return err

	case "list":
		List(fileName, writer)

	case "remove":
		id := args.id
		if id == "" {
			err = fmt.Errorf("%w", IdMissing)
		}
		Remove(id, fileName, writer)

	case "findById":
		id := args.id
		if id == "" {
			err = fmt.Errorf("%w", IdMissing)
		}
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
