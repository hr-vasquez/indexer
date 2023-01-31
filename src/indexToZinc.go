package main

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const INDEX_NAME = "email_index"
const JSON_FILE_NAME = "data.json"

// Need to be moved to a configuration file
const ZINCSEARCH_BULK_URL = "http://localhost:4080/api/_bulkv2"
const ZINCSEARCH_USER = "admin"
const ZINCSEARCH_PASSWORD = "admin"

func IndexSource(filePath string) *http.Response {
	logger.Println("Starting indexing to ZincSearch...")

	jsonPath := buildJsonFromSourceFiles(filePath)

	jsonFile, errOnOpen := os.Open(jsonPath)
	handleError(errOnOpen)

	client := http.Client{}
	request, errRequest := http.NewRequest("POST", ZINCSEARCH_BULK_URL, jsonFile)
	handleError(errRequest)
	request.SetBasicAuth(ZINCSEARCH_USER, ZINCSEARCH_PASSWORD)

	response, errResponse := client.Do(request)
	handleError(errResponse)

	return response
}

func buildJsonFromSourceFiles(sourcePath string) string {
	jsonFile, errCreatingFile := os.Create(JSON_FILE_NAME)
	handleError(errCreatingFile)

	writeContent(jsonFile, sourcePath, INDEX_NAME)

	errOnClose := jsonFile.Close()
	handleError(errOnClose)

	return JSON_FILE_NAME
}

func writeContent(jsonFile *os.File, sourcePath string, indexValue string) {
	var recordData []map[string]string

	allFilesFromPath := getFilesFromFolder(sourcePath)

	for i := 0; i < len(allFilesFromPath); i++ {
		contentMap := parseContentToJson(allFilesFromPath[i])
		recordData = append(recordData, contentMap)
	}

	bodyContent := BodyContent{Index: indexValue, Records: recordData}

	bodyJson, errParsing := json.Marshal(bodyContent)
	handleError(errParsing)

	_, errWriting := io.WriteString(jsonFile, string(bodyJson))
	handleError(errWriting)
}

func getFilesFromFolder(folderPath string) []string {
	folder, errOpeningFolder := os.Open(folderPath)
	handleError(errOpeningFolder)

	fileInfo, errGettingStats := folder.Stat()
	handleError(errGettingStats)

	if !fileInfo.IsDir() {
		return []string{folderPath}
	}

	var files []string

	errGettingChildren := filepath.Walk(folder.Name(),
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})
	handleError(errGettingChildren)

	return files
}

func parseContentToJson(path string) map[string]string {
	lines := getLinesFromFile(path)
	contentMap := buildMapFromLines(lines)
	return contentMap
}

func getLinesFromFile(path string) []string {
	file, errOnOpen := os.Open(path)
	handleError(errOnOpen)

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	errOnClose := file.Close()
	handleError(errOnClose)

	return lines
}

func buildMapFromLines(lines []string) map[string]string {
	dataMap := make(map[string]string)

	dataKeyList := []string{
		"Message-ID:",
		"Date:",
		"From:",
		"To:",
		"Subject:",
		"Cc:",
		"Mime-Version:",
		"Content-Type:",
		"Content-Transfer-Encoding:",
		"Bcc:",
		"X-From:",
		"X-To:",
		"X-cc:",
		"X-bcc:",
		"X-Folder:",
		"X-Origin:",
		"X-FileName:"}
	fillContentToMap(0, lines, 0, dataKeyList, dataMap)

	return dataMap
}

func fillContentToMap(
	lineIndex int,
	lines []string,
	keyIndex int,
	dataKeyList []string,
	dataMap map[string]string) {

	if lineIndex >= len(lines) {
		return
	}

	line := lines[lineIndex]

	if strings.TrimSpace(line) == "" {
		restOfLines := lines[lineIndex:]
		dataMap["Body"] = strings.Join(restOfLines, "\n")
		return
	}

	indexOfLineInKeys := findIndexOfValueInKeys(line, dataKeyList)
	if indexOfLineInKeys > -1 {
		keyIndex = indexOfLineInKeys

		keyItem := dataKeyList[indexOfLineInKeys]
		keyItemLength := len(keyItem)

		key := line[0 : keyItemLength-1] // For example: "SUBJECT:" -> "SUBJECT"
		value := line[keyItemLength:]
		dataMap[key] = value

	} else {
		// When you have many lines of a specific property, for example CC:
		currentKey := strings.TrimSuffix(dataKeyList[keyIndex], ":")
		dataMap[currentKey] = dataMap[currentKey] + " " + strings.TrimSpace(line)
	}

	lineIndex++
	fillContentToMap(lineIndex, lines, keyIndex, dataKeyList, dataMap)
}

func findIndexOfValueInKeys(value string, list []string) int {
	for i := 0; i < len(list); i++ {
		if strings.Contains(value, list[i]) {
			return i
		}
	}
	return -1
}
