package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const fileName = "test.json"
const filePermission = 0644

// Common validation tests
func TestOperationMissingError(t *testing.T) {
	var buffer bytes.Buffer

	expectedError := "-operation flag has to be specified"
	args := Arguments{
		Id:        "",
		Operation: "",
		Item:      "",
		FileName:  fileName,
	}
	err := Perform(args, &buffer)

	if err == nil {
		t.Error("Expect error when -operation flag is missing")
	}

	if err.Error() != expectedError {
		t.Errorf("Expect error to be '%s', but got '%s'", expectedError, err.Error())
	}
}

func TestWrongOperationError(t *testing.T) {
	var buffer bytes.Buffer
	args := Arguments{
		Id:        "",
		Operation: "abcd",
		Item:      "",
		FileName:  fileName,
	}
	expectedError := "Operation abcd not allowed!"

	err := Perform(args, &buffer)

	if err == nil {
		t.Error("Expect error when wrong -operation passed")
	}

	if err.Error() != expectedError {
		t.Errorf("Expect error to be '%s', but got '%s'", expectedError, err.Error())
	}
}

func TestFileNameMissingError(t *testing.T) {
	var buffer bytes.Buffer
	args := Arguments{
		Id:        "",
		Operation: "list",
		Item:      "",
		FileName:  "",
	}
	expectedError := "-fileName flag has to be specified"

	err := Perform(args, &buffer)

	if err == nil {
		t.Error("Expect error when -fileName flag is missing")
	}

	if err.Error() != expectedError {
		t.Errorf("Expect error to be '%s', but got '%s'", expectedError, err.Error())
	}
}

// List operation tests
func TestListOperation(t *testing.T) {
	args := Arguments{
		Id:        "",
		Operation: "list",
		Item:      "",
		FileName:  fileName,
	}
	var buffer bytes.Buffer

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, filePermission)
	defer os.Remove(fileName)
	if err != nil {
		t.Error(err)
	}

	existingItems := "{\"users\": [{\"id\": \"1\", \"email\": \"test@test.com\", \"age\": 34},{\"id\": \"2\", \"email\": \"tes2@test.com\", \"age\": 32}]}"

	file.Write([]byte(existingItems))
	file.Close()
	err = Perform(args, &buffer)
	if err != nil {
		t.Error(err)
	}

	file, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, filePermission)
	if err != nil {
		t.Error(err)
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		t.Error(err)
	}

	result := fmt.Sprintf("{\"users\": %s}", strings.TrimRight(buffer.String(), "\n"))
	if result != existingItems {
		t.Errorf("Expect output to equal %s, but got %s", existingItems, result)
	}
	if string(bytes) != existingItems {
		t.Errorf("Expect file content to equal %s, but got %s", existingItems, string(bytes))
	}
}

//Adding operation tests
//func TestAddingOperationMissingItem(t *testing.T) {
//	var buffer bytes.Buffer
//	args := Arguments{
//		Id:        "",
//		Operation: "add",
//		Item:      "",
//		FileName:  fileName,
//	}
//	expectedError := "-item flag has to be specified"
//	defer os.Remove(fileName)
//
//	err := Perform(args, &buffer)
//
//	if err == nil {
//		t.Error("Expect error when -item flag is missing")
//	}
//
//	if err.Error() != expectedError {
//		t.Errorf("Expect error to be '%s', but got '%s'", expectedError, err.Error())
//	}
//}

func TestAddingOperationSameID(t *testing.T) {
	var buffer bytes.Buffer

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, filePermission)
	defer os.Remove(fileName)

	if err != nil {
		t.Error(err)
	}

	existingItem := "{\"users\": [{\"id\": \"1\", \"email\": \"test@test.com\", \"age\": 34}]}"

	file.Write([]byte(existingItem))
	file.Close()

	item := "{\"id\": \"1\", \"email\": \"test@test.com\", \"age\": 34}"
	args := Arguments{
		Id:        "",
		Operation: "add",
		Item:      item,
		ItemUsr:   User{"1", "test@test.com", 34},
		FileName:  fileName,
	}
	expectedOutput := "Item with id 1 already exists"

	err = Perform(args, &buffer)
	if err != nil {
		t.Error(err)
	}

	resultOutput := buffer.String()

	if resultOutput != expectedOutput {
		t.Errorf("Expect error to be '%s', but got '%s'", expectedOutput, resultOutput)
	}
}

func TestAddingOperation(t *testing.T) {
	var buffer bytes.Buffer

	expectedFileContent := "{\"users\":[{\"Id\":\"1\",\"Email\":\"test@test.com\",\"Age\":34}]}"
	itemToAdd := "{\"id\": \"1\", \"email\": \"test@test.com\", \"age\": 34}"
	args := Arguments{
		Id:        "",
		Operation: "add",
		Item:      itemToAdd,
		ItemUsr:   User{"1", "test@test.com", 34},
		FileName:  fileName,
	}
	//defer os.Remove(fileName)
	err := Perform(args, &buffer)
	if err != nil {
		t.Error(err)
	}

	file, err := os.OpenFile(fileName, os.O_RDONLY, filePermission)
	defer file.Close()

	if err != nil {
		t.Error(err)
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		t.Error(err)
	}

	if string(bytes) != expectedFileContent {
		t.Errorf("Expect file content to be %s, but got %s", expectedFileContent, bytes)
	}
}

// FindByID operation tests
//func TestFindByIdOperationMissingID(t *testing.T) {
//	var buffer bytes.Buffer
//	args := Arguments{
//		Id:        "",
//		Operation: "findById",
//		Item:      "",
//		FileName:  fileName,
//	}
//	expectedError := "-id flag has to be specified"
//
//	err := Perform(args, &buffer)
//
//	if err == nil {
//		t.Error("Expect error when -id flag is missing")
//	}
//
//	if err.Error() != expectedError {
//		t.Errorf("Expect error to be '%s', but got '%s'", expectedError, err.Error())
//	}
//}

//func TestFindByIdOperation(t *testing.T) {
//	var buffer bytes.Buffer
//
//	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, filePermission)
//	defer os.Remove(fileName)
//
//	if err != nil {
//		t.Error(err)
//	}
//
//	existingItems := "[{\"id\":\"1\",\"email\":\"test@test.com\",\"age\":34},{\"id\":\"2\",\"email\":\"test2@test.com\",\"age\":31}]"
//
//	file.Write([]byte(existingItems))
//	file.Close()
//
//	expectedOutput := "{\"id\":\"2\",\"email\":\"test2@test.com\",\"age\":31}"
//	args := Arguments{
//		Id:        "2",
//		Operation: "findById",
//		Item:      "",
//		FileName:  fileName,
//	}
//	err = Perform(args, &buffer)
//
//	if err != nil {
//		t.Error(err)
//	}
//
//	resultString := buffer.String()
//
//	if resultString != expectedOutput {
//		t.Errorf("Expect output to be '%s', but got '%s'", expectedOutput, resultString)
//	}
//}

//func TestFindByIdOperationWrongID(t *testing.T) {
//	var buffer bytes.Buffer
//
//	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, filePermission)
//	defer os.Remove(fileName)
//
//	if err != nil {
//		t.Error(err)
//	}
//
//	existingItems := "[{\"id\":\"1\",\"email\":\"test@test.com\",\"age\":34},{\"id\":\"2\",\"email\":\"test2@test.com\",\"age\":31}]"
//
//	file.Write([]byte(existingItems))
//	file.Close()
//
//	expectedOutput := ""
//	args := Arguments{
//		Id:        "3",
//		Operation: "findById",
//		Item:      "",
//		FileName:  fileName,
//	}
//	err = Perform(args, &buffer)
//
//	if err != nil {
//		t.Error(err)
//	}
//
//	resultString := buffer.String()
//
//	if resultString != expectedOutput {
//		t.Errorf("Expect output to be '%s', but got '%s'", expectedOutput, resultString)
//	}
//}

// Removing operations tests

//func TestRemovingOperationMissingID(t *testing.T) {
//	var buffer bytes.Buffer
//	args := Arguments{
//		Id:        "",
//		Operation: "remove",
//		Item:      "",
//		FileName:  fileName,
//	}
//
//	expectedError := "-id flag has to be specified"
//
//	err := Perform(args, &buffer)
//
//	if err == nil {
//		t.Error("Error has to be shown when -id flag is missing")
//	}
//
//	if err.Error() != expectedError {
//		t.Errorf("Expect error to be '%s', but got '%s'", expectedError, err.Error())
//	}
//}

//func TestRemovingOperationWrongID(t *testing.T) {
//	var buffer bytes.Buffer
//
//	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, filePermission)
//	defer os.Remove(fileName)
//
//	if err != nil {
//		t.Error(err)
//	}
//
//	existingItems := "[{\"id\":\"1\",\"email\":\"test@test.com\",\"age\":34}]"
//
//	file.Write([]byte(existingItems))
//	file.Close()
//
//	expectedOutput := "Item with id 2 not found"
//	args := Arguments{
//		Id:        "2",
//		Operation: "remove",
//		Item:      "",
//		FileName:  fileName,
//	}
//	err = Perform(args, &buffer)
//
//	if err != nil {
//		t.Error(err)
//	}
//
//	resultOutput := buffer.String()
//
//	if resultOutput != expectedOutput {
//		t.Errorf("Expect output to be '%s', but got '%s'", expectedOutput, resultOutput)
//	}
//}

//func TestRemovingOperation(t *testing.T) {
//	var buffer bytes.Buffer
//
//	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, filePermission)
//	defer os.Remove(fileName)
//
//	if err != nil {
//		t.Error(err)
//	}
//
//	existingItems := "[{\"id\":\"1\",\"email\":\"test@test.com\",\"age\":34},{\"id\":\"2\",\"email\":\"test2@test.com\",\"age\":31}]"
//
//	file.Write([]byte(existingItems))
//	file.Close()
//	expectedFileContent := "[{\"id\":\"2\",\"email\":\"test2@test.com\",\"age\":31}]"
//	args := Arguments{
//		Id:        "1",
//		Operation: "remove",
//		Item:      "",
//		FileName:  fileName,
//	}
//
//	err = Perform(args, &buffer)
//
//	if err != nil {
//		t.Error(err)
//	}
//
//	file, err = os.OpenFile(fileName, os.O_RDONLY, filePermission)
//	defer file.Close()
//
//	if err != nil {
//		t.Error(err)
//	}
//
//	bytes, err := ioutil.ReadAll(file)
//
//	if string(bytes) != expectedFileContent {
//		t.Errorf("Expect file content to be '%s', but got '%s'", expectedFileContent, bytes)
//	}
//}
