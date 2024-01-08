package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
)

const Version = "1.0.1"

type (
	Logger interface {
		Fatal(string, ...interface{})
		Error(string, ...interface{})
		Warn(string, ...interface{})
		Info(string, ...interface{})
		Debug(string, ...interface{})
		Trace(string, ...interface{})
	}

	Driver struct {
		mutex   sync.Mutex
		mutexes map[string]*sync.Mutex
		dir     string
		log     Logger
	}
)

type Options struct {
	Logger Logger
}

func New(dir string, options *Options) (*Driver, error) {
	dir = filepath.Clean(dir)

	if options == nil {
		options = &Options{}
	}

	if options.Logger == nil {
		options.Logger = lumber.NewConsoleLogger(lumber.INFO)
	}

	driver := &Driver{
		mutex:   sync.Mutex{},
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
		log:     options.Logger,
	}

	if _, err := stat(dir); err == nil {
		options.Logger.Debug("Using existing database at %s", dir)
		return driver, nil
	}
	options.Logger.Debug("Creating new database at %s", dir)
	return driver, os.MkdirAll(dir, 0755)

}

func (driver *Driver) Write(collection, resource string, value interface{}) error {
	if collection == "" {
		return fmt.Errorf("collection cannot be empty")
	}

	if resource == "" {
		return fmt.Errorf("resource cannot be empty")
	}

	mutex := driver.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(driver.dir, collection)
	fnlPath := filepath.Join(dir, resource+".json")
	tmpPath := fnlPath + ".tmp"

	if err := os.Mkdir(dir, 0755); err != nil && !os.IsExist(err) {
		return err
	}

	b, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}

	b = append(b, byte('\n'))

	if err := ioutil.WriteFile(tmpPath, b, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, fnlPath)
}

func (driver *Driver) Read(collection, resource string, value interface{}) (string, error) {
	if collection == "" {
		return "", fmt.Errorf("collection cannot be empty")
	}

	if resource == "" {
		return "", fmt.Errorf("resource cannot be empty")
	}

	record := filepath.Join(driver.dir, collection, resource)

	if _, err := stat(record); err != nil {
		return "", err
	}

	b, err := ioutil.ReadFile(record + ".json")
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(b, value); err != nil {
		return "", err
	}

	return string(b), nil
}

func (driver *Driver) ReadAll(collection string) ([]string, error) {
	if collection == "" {
		return nil, fmt.Errorf("collection cannot be empty")
	}

	dir := filepath.Join(driver.dir, collection)
	if _, err := stat(dir); err != nil {
		return nil, err
	}

	files, _ := ioutil.ReadDir(dir)
	// if err != nil {
	// 	return nil, err
	// }

	fmt.Println("Dir: ", files)
	var records []string

	for _, file := range files {

		b, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		records = append(records, string(b))

	}

	return records, nil
}

func (driver *Driver) Delete(collection, resource string) error {
	if collection == "" {
		return fmt.Errorf("collection cannot be empty")
	}

	if resource == "" {
		return fmt.Errorf("resource cannot be empty")
	}

	mutex := driver.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	record := filepath.Join(driver.dir, collection, resource)

	switch fi, err := stat(record); {
	case err != nil, fi == nil:
		return fmt.Errorf("cannot delete resource: %s", err)
	case fi.Mode().IsDir():
		return os.RemoveAll(record)
	case fi.Mode().IsRegular():
		return os.RemoveAll(record + ".json")
	}
	return nil
}

func stat(path string) (fi os.FileInfo, err error) {
	if fi, err = os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path + ".json")
	}
	return fi, err
}

func (driver *Driver) getOrCreateMutex(collection string) *sync.Mutex {
	driver.mutex.Lock()
	defer driver.mutex.Unlock()

	mutex, ok := driver.mutexes[collection]
	if !ok {
		mutex = &sync.Mutex{}
		driver.mutexes[collection] = mutex
	}
	return mutex
}

type User struct {
	Name    string
	Age     json.Number
	Contact string
	Comapny string
	Address Address
}

type Address struct {
	Street  string
	City    string
	State   string
	Pincode json.Number
}

func main() {
	dir := "./"
	db, err := New(dir, nil)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	employess := []User{
		{
			Name:    "John",
			Age:     "30",
			Contact: "1234567890",
			Comapny: "ABC",
			Address: Address{
				Street:  "Street 1",
				City:    "City 1",
				State:   "State 1",
				Pincode: "123456",
			},
		},
		{
			Name:    "Doe",
			Age:     "40",
			Contact: "1234567890",
			Comapny: "ABC",
			Address: Address{
				Street:  "Street 2",
				City:    "City 2",
				State:   "State 2",
				Pincode: "123456",
			},
		},
		{
			Name:    "Smith",
			Age:     "50",
			Contact: "1234567890",
			Comapny: "ABC",
			Address: Address{
				Street:  "Street 3",
				City:    "City 3",
				State:   "State 3",
				Pincode: "123456",
			},
		},
		{
			Name:    "Alex",
			Age:     "60",
			Contact: "1234567890",
			Comapny: "ABC",
			Address: Address{
				Street:  "Street 4",
				City:    "City 4",
				State:   "State 4",
				Pincode: "123456",
			},
		},
		{
			Name:    "Amit",
			Age:     "30",
			Contact: "1234567890",
			Comapny: "ABC",
			Address: Address{
				Street:  "Street 5",
				City:    "City 5",
				State:   "State 5",
				Pincode: "123456",
			},
		},
	}

	for _, employee := range employess {
		db.Write("users", employee.Name, employee)
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	records, err := db.ReadAll("users")
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("Records: ", records)

	allUsers := []User{}
	for _, record := range records {
		user := User{}
		err := json.Unmarshal([]byte(record), &user)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		allUsers = append(allUsers, user)
	}

	fmt.Println("All Users: ", allUsers)

	user, err := db.Read("users", "John", &User{})
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("User: ", user)

	err = db.Delete("users", "John")
	if err != nil {
		fmt.Println("Error: ", err)
	}

}
