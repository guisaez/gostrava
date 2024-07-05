package gostrava

import (
	"fmt"
	"os"
	"path/filepath"
)

func WriteJsonFile(file_name string, json_data []byte) {
	folderPath := "sample_response"
	
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating folder", err)
	}

	file, err := os.Create(filepath.Join(folderPath, file_name))
	if err != nil {
		fmt.Println("Error creating file: ", err)
	}
	defer file.Close()

	_, err = file.Write(json_data)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}
