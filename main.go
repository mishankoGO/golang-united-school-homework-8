package main

import (
	"fmt"
	"io"
	"os"
)

// User struct for storing user data
type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// Arguments type to store operation, item, id and fileName
type Arguments map[string]string

// Perform function that performs all the application logic
func Perform(args map[string]string, writer io.Writer) (err error) {

	op := args["operation"]
	fileName := args["fileName"]

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
		if len(args["item"]) != 0 {
			itemUsr := ItemParser(args)
			Add(itemUsr, fileName, writer)
		} else {
			err = fmt.Errorf("%w", ItemMissing)
			return err
		}

	case "list":
		List(fileName, writer)

	case "remove":
		id := args["id"]
		if id == "" {
			err = fmt.Errorf("%w", IdMissing)
		}
		Remove(id, fileName, writer)

	case "findById":
		id := args["id"]
		if id == "" {
			err = fmt.Errorf("%w", IdMissing)
		}
		FindById(id, fileName, writer)

	// if no operation connects - throw an error
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
