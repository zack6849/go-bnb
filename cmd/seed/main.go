package main

import (
	"GoBNB/internal/bootstrap"
	"GoBNB/internal/configuration"
	"GoBNB/internal/domain"
	"GoBNB/internal/importer"
	"GoBNB/internal/query_cache"
	"fmt"
	"io/fs"
	"os"

	"gorm.io/gorm"
)

func main() {
	_ = bootstrap.Initialize()
	configuration.Load()
	db, err := configuration.GetDatabaseConfiguration().Open()
	if err != nil {
		fmt.Printf("failed to open db connection: %s", err.Error())
		return
	}
	seedTypes := []string{"listings"}
	for _, seedType := range seedTypes {
		err := importSeedType(seedType, db)
		if err != nil {
			fmt.Printf("Failed to import seed type %s:\n\t%s", seedType, err.Error())
		}
	}
}

func importSeedType(seedType string, db *gorm.DB) error {
	path := "seed_data/" + seedType + ".csv"
	cache := query_cache.NewCache()
	file, err := getSeedReader(path)
	if err != nil {
		return err
	}
	//make sure we close this file handle
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			//nothing we can really do, but i guess we should at least tell people
			fmt.Printf("failed to close file %s: %s", path, err.Error())
		}
	}(file)
	total := 0
	//print total insert count at the end
	defer func() {
		fmt.Printf("Inserted %d records total\r\n", total)
	}()
	//stream the file into a map of strings and strings
	return importer.StreamHeaderAwareImport[map[string]string](
		file,
		func(header []string, record []string) (map[string]string, error) {
			return importer.ArrayCombine(header, record)
		},
		10,
		func(batch []map[string]string) error {
			total += len(batch)
			insertBatch := make([]domain.Listing, 0, len(batch))
			index := total
			for _, item := range batch {
				listing, err := CreateListingRecord(item, cache, db)
				listingId := item["id"]
				index++
				if err != nil {
					return fmt.Errorf("listing#%d (%s) import failed:\n\t\t%w", index, listingId, err)
				}
				insertBatch = append(insertBatch, listing)
			}
			if res := db.CreateInBatches(insertBatch, 10); res.Error != nil {
				return fmt.Errorf("failed to insert listing batch: %w", res.Error)
			}
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
