package helper

import (
	"fmt"
	"log"
	"os"
)

func InitLogging(buildtime string) {
	file, err := openLogFile(fmt.Sprintf("./mastery-helper-%s.log", buildtime))
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

}

func openLogFile(path string) (*os.File, error) {
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return logFile, nil
}
