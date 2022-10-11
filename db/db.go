package db

import (
	"errors"
	"fmt"
	"github.com/alidevjimmy/go-raft/fileutils"
	"strconv"
	"strings"
)

const (
	KeyNotExistsErr = "key not exists"
)
const (
	GetOperation    = "GET"
	DeleteOperation = "DELETE"
	SetOperation    = "SET"
)

// in memory log
type Database struct {
	db map[string]int
}

func NewDatabase() *Database {
	kv := make(map[string]int)
	return &Database{
		db: kv,
	}
}
func (d *Database) getKey(key string) (int, error) {
	val, ok := d.db[key]
	if !ok {
		return -1, errors.New(KeyNotExistsErr)
	}
	return val, nil
}

// adds if key not exists
// updates if key exists
func (d *Database) setKey(key string, val int) error {
	d.db[key] = val
	return nil
}

func (d *Database) deleteKey(key string) error {
	_, exists := d.db[key]
	if !exists {
		return errors.New(KeyNotExistsErr)
	}
	delete(d.db, key)
	return nil
}

func (d *Database) PersistLogCommand(command, serverName string) error {
	fileName := fmt.Sprintf("%s.log", serverName)
	err := fileutils.CreateFileIfNotExists(fileName)
	if err != nil {
		return err
	}
	err = fileutils.WriteToFile(fileName, fmt.Sprintf("%s:%s\n", serverName, command))
	if err != nil {
		return err
	}
	return nil
}

// GET key
// DELETE key
// SET key value
func (d *Database) ValidateCommand(command string) error {
	splits := strings.Split(command, " ")
	op := splits[0]
	if op == GetOperation || op == DeleteOperation {
		if len(splits) != 2 {
			return errors.New("GET & DELETE operations needs key")
		}
	} else if op == SetOperation {
		if len(splits) != 3 {
			return errors.New("SET operation need key and value")
		}
		_, err := strconv.Atoi(splits[2])
		if err != nil {
			return errors.New("value should be integer")
		}
	} else {
		return errors.New("invalid command")
	}
	return nil
}

func (d *Database) PerformCommand(command string) string {
	splits := strings.Split(command, " ")
	op := splits[0]
	key := splits[1]
	res := ""
	switch op {
	case GetOperation:
		val, err := d.getKey(key)
		if err != nil {
			res = err.Error()
		} else {
			res = fmt.Sprintf("value of key (%s) is: %d\n", key, val)
		}
	case DeleteOperation:
		err := d.deleteKey(key)
		if err != nil {
			res = err.Error()
		} else {
			res = fmt.Sprintf("key (%s) deleted", key)
		}
	case SetOperation:
		val := splits[2]
		intVal, _ := strconv.Atoi(val)
		if err := d.setKey(key, intVal); err != nil {
			res = err.Error()
		} else {
			res = "value stored in storage"
		}
	}
	return res
}
