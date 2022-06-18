package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// helper function for "add" operation
func Add(item User, fileName string, writer io.Writer) {
	// open/create file
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, Permission)
	if err != nil {
		fmt.Println("error creating file")
		panic(err)
	}
	defer f.Close()

	// read all the contents
	buf, _ := ioutil.ReadAll(f)

	// fill in the User json
	var users []User
	_ = json.Unmarshal(buf, &users)

	// check whether id exists or not
	for _, elem := range users {
		if elem.Id == item.Id {
			writer.Write([]byte(IdExists(elem.Id)))
		}
	}
	if item.Age != 0 {
		users = append(users, item)
	}

	// write updated users
	updUsers, _ := json.Marshal(&users)
	os.WriteFile(fileName, updUsers, Permission)
}

// helper function for "list" operation
func List(fileName string, writer io.Writer) {
	// open/create file
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, Permission)
	if err != nil {
		fmt.Println("error creating file")
		panic(err)
	}
	defer f.Close()

	// read all the contents and write them
	buf, _ := ioutil.ReadAll(f)
	if len(buf) != 0 {
		writer.Write(buf)
	}
}

// helper function for "remove" operation
func Remove(id string, fileName string, writer io.Writer) {
	// read file
	f, _ := os.ReadFile(fileName)

	// read json
	var users, newUsers []User
	_ = json.Unmarshal(f, &users)
	for i := 0; i < len(users); i++ {
		// if id in json - then delete
		if users[i].Id != id {
			newUsers = append(newUsers, users[i])
		}
	}

	// write updated json
	updUsers, _ := json.Marshal(&newUsers)

	if len(newUsers) == len(users) {
		writer.Write([]byte(IdNotFound(id)))
	}
	os.WriteFile(fileName, updUsers, Permission)
}

// helper function for "findById" operation
func FindById(id string, fileName string, writer io.Writer) {
	// read file
	f, _ := os.ReadFile(fileName)

	//read json
	var users []User
	_ = json.Unmarshal(f, &users)
	for _, elem := range users {
		// if found id - then write
		if elem.Id == id {
			res := fmt.Sprintf("{\"id\":\"%s\",\"email\":\"%s\",\"age\":%d}", elem.Id, elem.Email, elem.Age)
			writer.Write([]byte(res))
			return
		}
	}
	writer.Write([]byte(""))
}

// helper function to parse item string
func ItemParser(a Arguments) User {
	item := a["item"]
	u := User{}

	// remove unnecessary characters
	item = strings.Trim(item, "[{}]")
	item = strings.Replace(item, " ", "", -1)
	itemLst := strings.Split(item, ",")

	// fill User struct
	for _, elem := range itemLst {
		kVal := stringSplitter(elem)
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

// helper function for string splitting
func stringSplitter(s string) []string {
	return strings.Split(s, ":")
}

// helper function to fill up the Arguments map
func parseArgs() Arguments {
	op := flag.String("operation", "", "operation to perform on users list")
	item := flag.String("item", "", "item to add")
	Id := flag.String("id", "", "user id")
	fileName := flag.String("fileName", "", "name of the file with users list")
	flag.Parse()

	var args = make(Arguments)

	args["operation"] = *op
	args["id"] = *Id
	if len(*item) != 0 {
		s := strings.Trim(*item, "{}")
		item := strings.Split(s, ",")
		id := strings.Trim(stringSplitter(item[0])[1], " ")
		email := strings.Trim(stringSplitter(item[1])[1], " ")
		age, err := strconv.Atoi(strings.Trim(stringSplitter(item[2])[1], " "))
		if err != nil {
			fmt.Println("error converting age to int")
		}
		args["item"] = fmt.Sprintf("[{\"id\" :\"%s\",\"email\" :\"%s\",\"age\" :%d}]", id, email, age)

	}

	args["fileName"] = *fileName

	return args
}
