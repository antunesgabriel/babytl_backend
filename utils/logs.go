package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func RegisterLog(indentify string, errMessage string) {
	date := time.Now().Format("02_01_2006")

	fileName := date + ".txt"

	file, err := os.OpenFile(filepath.Join("logs", fileName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)

	if err != nil {
		file.Close()
		fmt.Println("Erro", err)
		return
	}

	defer file.Close()

	timestamps := time.Now().Format("02/01/2006 15:04:05")
	message := fmt.Sprintf("[%s] : [%s] - Time: %s", indentify, errMessage, timestamps)

	if _, err := file.WriteString(message); err != nil {
		fmt.Println(err)
		return
	}

}
