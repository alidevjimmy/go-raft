package fileutils

import "os"

func CreateFileIfNotExists(name string) (err error) {
	_, err = os.Stat(name)
	if os.IsNotExist(err) {
		f, err := os.Create(name)
		defer func() {
			err = f.Close()
		}()
		if err != nil {
			return err
		}
	} else {
		return
	}
	return
}

func WriteToFile(name, message string) error {
	f, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = f.WriteString(message)
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}
