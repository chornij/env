package env

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// https://github.com/joho/godotenv/blob/master/godotenv.go

func TryLoad(envFile string) {
	if _, err := os.Stat(envFile); err == nil || os.IsExist(err) {
		if err := Load(envFile); err != nil {
			log.Println("load env >", err)
		}
	}
}

func Get(name, defaultValue string) string {
	value := os.Getenv(name)

	if value == "" {
		return defaultValue
	}

	return value
}

func GetInt(name string, defaultValue int) (result int) {
	value := os.Getenv(name)

	if value != "" {
		result, _ = strconv.Atoi(value)
	} else {
		result = defaultValue
	}

	return
}

func GetInt64(name string, defaultValue int64) (result int64) {
	value := os.Getenv(name)

	if value != "" {
		result, _ = strconv.ParseInt(name, 10, 64)
	} else {
		result = defaultValue
	}

	return
}

func Load(filename string) error {
	envMap, err := readFile(filename)

	if err != nil {
		return err
	}

	currentEnv := map[string]bool{}
	rawEnv := os.Environ()

	for _, rawEnvLine := range rawEnv {
		key := strings.Split(rawEnvLine, "=")[0]
		currentEnv[key] = true
	}

	for key, value := range envMap {
		if !currentEnv[key] {
			if err := os.Setenv(key, value); err != nil {
				fmt.Println(err)

				return err
			}
		}
	}

	return nil
}

func readFile(filename string) (envMap map[string]string, err error) {
	file, err := os.Open(filename)

	if err != nil {
		return
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()

	return parse(file)
}

// parse reads an env file from io.Reader, returning a map of keys and values.
func parse(r io.Reader) (envMap map[string]string, err error) {
	envMap = make(map[string]string)

	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return
	}

	for _, fullLine := range lines {
		var key, value string
		key, value, err = parseLine(fullLine)

		if err != nil {
			return
		}

		envMap[key] = value
	}

	return
}

func parseLine(line string) (key string, value string, err error) {
	if len(line) == 0 {
		err = errors.New("zero length string")

		return
	}

	info := strings.SplitN(line, "=", 2)

	if len(info) != 2 {
		err = errors.New("can't separate key from value: " + line)

		return
	}

	key = strings.Trim(info[0], " ")
	value = strings.Trim(info[1], " ")

	return
}
