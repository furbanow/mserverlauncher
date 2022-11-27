package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ParseServerPropertiesFile(srvPropsPath string) (map[string]string, error) {
	file, err := os.Open(srvPropsPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	return ParseServerProperties(scanner), nil
}

func ParseServerPropertiesString(content string) map[string]string {
	reader := strings.NewReader(content)
	scanner := bufio.NewScanner(reader)
	return ParseServerProperties(scanner)
}

func ParseServerProperties(scanner *bufio.Scanner) map[string]string {

	props := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") {
			continue
		}

		for strings.Contains(line, "#") {
			line = strings.Split(line, "#")[0]
		}

		parts := strings.Split(line, "=")
		if len(parts) == 1 {
			props[parts[0]] = ""
		} else if len(parts) == 2 {
			props[parts[0]] = parts[1]
		}
	}
	return props
}

func SaveServerPropertiesFile(props map[string]string, path string) error {

	file, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for key, elem := range props {
		_, err := writer.WriteString(fmt.Sprintf("%s=%s\n", key, elem))
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	writer.Flush()
	return nil
}
