package main

import (
	"GoBNB/internal/configuration"
	"GoBNB/internal/importer"
	"fmt"
	"io/fs"
	"os"
)

func main() {
	configuration.Load()
	seedTypes := []string{"listings"}
	for _, seedType := range seedTypes {
		err := importSeedType(seedType)
		if err != nil {
			fmt.Printf("Failed to import seed type %s: %s", seedType, err.Error())
		}
	}
}

func ArrayCombine(keys []string, values []string) (map[string]string, error) {
	data := make(map[string]string)
	if len(keys) != len(values) {
		return nil, fmt.Errorf("row does not have the same number of values as the header row")
	}
	for i := range len(keys) {
		header := keys[i]
		value := values[i]
		data[header] = value
	}
	return data, nil
}

func importSeedType(seedType string) error {
	path := "seed_data/" + seedType + ".csv"
	file, err := getSeedReader(path)
	if err != nil {
		return err
	}
	//make sure we close this file handle
	defer file.Close()
	total := 0
	defer func() {
		fmt.Printf("Inserted %d records total\r\n", total)
	}()

	//stream the file into a map of strings and strings
	return importer.StreamHeaderAwareImport[map[string]string](
		file,
		func(header []string, record []string) (map[string]string, error) {
			return ArrayCombine(header, record)
		},
		10,
		func(batch []map[string]string) error {
			total += len(batch)
			fmt.Printf("Inserting %d listings!\r\n", len(batch))
			return nil
		},
	)
}

func getSeedReader(path string) (*os.File, error) {
	if fs.ValidPath(path) {
		file, err := os.Open(path)
		if err != nil {
			return file, err
		}
		return file, nil
	}
	return nil, fmt.Errorf("path %s is not valid", path)
}
