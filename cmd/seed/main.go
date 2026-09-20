package main

import (
	"context"
	"flag"
	"fmt"
	"gobnb/internal/bootstrap"
	"gobnb/internal/cache"
	"gobnb/internal/configuration"
	"gobnb/internal/domain"
	"gobnb/internal/errs"
	"gobnb/internal/importer"
	"gobnb/internal/services/geoapify"
	"gobnb/internal/services/geocode"
	"io"
	"io/fs"
	"log"
	"os"
	"time"

	"github.com/mmcloughlin/geohash"
	"gorm.io/gorm"
)

func main() {
	ctx := bootstrap.Initialize()
	configuration.Load()

	var skipGeo bool
	var filename string

	flag.StringVar(&filename, "file", "listings.csv", "filename to import (relative to seed_data dir)")
	flag.BoolVar(&skipGeo, "skipgeo", false, "if we should skip geocoding on import")
	flag.Parse()

	handle, err := resolveFileByName(filename)
	if err != nil {
		fmt.Printf("failed to resolve import filename %s: %s", filename, err.Error())
	}

	//make sure we close the file handle
	defer func(handle *os.File) {
		err := handle.Close()
		if err != nil {
			fmt.Printf("failed to close seed file %s: %s", handle.Name(), err.Error())
		}
	}(handle)

	db, err := configuration.GetDatabaseConfiguration().Open()

	log.Printf("Seeding started (%s)", handle.Name())
	if err != nil {
		fmt.Printf("failed to open db connection: %s", err.Error())
		return
	}

	start := time.Now()
	log.Println("Resolving dependent records (property types, amenities, etc)")
	if err := createDependentRecords(handle, db); err != nil {
		fmt.Printf("failed to create dependent records from file %s: %s", handle.Name(), err.Error())
		return
	}
	log.Println("Importing listings in bulk")
	if err := importListingRecords(handle, db); err != nil {
		fmt.Printf("failed to import listing records from file %s: %s\n", handle.Name(), err.Error())
		return
	}
	if !skipGeo {
		log.Println("Reverse-geocoding listings to prime cities cache...")
		if err := reverseGeocodeListings(handle, db, ctx); err != nil {
			fmt.Printf("failed to geocode listing records: %s", err.Error())
			return
		}
		log.Println("Geocoding complete")
	}

	end := time.Now()
	elapsed := end.Sub(start).Abs().String()
	log.Println("import complete. finished in " + elapsed)
}

func createDependentRecords(file *os.File, db *gorm.DB) error {
	err := rewindFileHandle(file)
	if err != nil {
		return err
	}
	const recordChunkSize = 2000
	c := cache.NewCache()
	err = importer.StreamHeaderAwareImport[KeyValueRecord](
		file,
		func(header []string, record []string) (KeyValueRecord, error) {
			//combine the header with the record
			return importer.ArrayCombine(header, record)
		},
		recordChunkSize,
		func(batch []KeyValueRecord) error {
			errBag := errs.CreateEchoChamber()
			err := db.Transaction(func(db *gorm.DB) error {
				for _, record := range batch {
					//parse the listing
					listing, err := ExtractListingFromRecord(record)
					errBag.PushError(err)

					//resolve the sub-relations
					roomType, err := resolveRoomType(c, record["room_type"], db)
					errBag.PushError(err)
					listing.RoomType = roomType

					propertyType, err := resolvePropertyType(c, record["property_type"], db)
					listing.PropertyType = propertyType
					errBag.PushError(err)

					amenityNames, err := ExtractAmenityNamesFromRecord(record)
					if err == nil {
						resolvedAmenities := resolveAmenities(c, amenityNames, db)
						listing.Amenities = resolvedAmenities
					}

					host, err := ExtractHostFromRecord(record)
					errBag.PushError(err)
					if err == nil {
						resolvedHost, err := resolveHost(c, host, db)
						errBag.PushError(err)
						listing.Hosts[0] = resolvedHost
					}
				}
				return db.Error
			})
			errBag.PushError(err)
			return errBag.Summary()
		},
	)
	return err
}

func importListingRecords(file *os.File, db *gorm.DB) error {
	const dbChunkSize = 250
	const recordChunkSize = 2000
	err := rewindFileHandle(file)
	if err != nil {
		return err
	}
	c := cache.NewCache()
	return importer.StreamHeaderAwareImport[KeyValueRecord](
		file,
		func(header []string, record []string) (KeyValueRecord, error) {
			//combine the header with the record
			return importer.ArrayCombine(header, record)
		},
		recordChunkSize,
		func(batch []KeyValueRecord) error {
			var bulk []domain.Listing
			for _, record := range batch {
				errBag := errs.CreateEchoChamber()
				//parse the listing
				listing, err := ExtractListingFromRecord(record)
				errBag.PushError(err)

				//resolve the sub-relations
				roomType, err := resolveRoomType(c, listing.RoomType.Name, db)
				errBag.PushError(err)
				listing.RoomType = roomType

				propertyType, err := resolvePropertyType(c, listing.PropertyType.Name, db)
				listing.PropertyType = propertyType
				errBag.PushError(err)

				host, err := resolveHost(c, listing.Hosts[0], db)
				errBag.PushError(err)
				listing.Hosts[0] = host

				amenityNames, err := ExtractAmenityNamesFromRecord(record)
				if err == nil {
					resolvedAmenities := resolveAmenities(c, amenityNames, db)
					listing.Amenities = resolvedAmenities
				}

				bulk = append(bulk, listing)
			}
			tx := db.Model(domain.Listing{}).CreateInBatches(bulk, dbChunkSize)
			return tx.Error
		},
	)
}

// geohash the location to normalize the locations before lookup and prevent a bunch of lookups that are relatively close to each other
func reverseGeocodeListings(file *os.File, db *gorm.DB, ctx context.Context) error {
	//rewind reader pointer to the start of the file
	_, err := file.Seek(0, 0)
	if err != nil {
		return err
	}
	queue := make(map[string]bool)

	//parse the listing and append a geohash for each listing to our dataset
	err = importer.StreamHeaderAwareImport[KeyValueRecord](
		file,
		func(header []string, record []string) (KeyValueRecord, error) {
			//combine the header with the record
			return importer.ArrayCombine(header, record)
		},
		50,
		func(batch []KeyValueRecord) error {
			for _, record := range batch {
				//parse the listing
				listing, err := ExtractListingFromRecord(record)
				if err != nil {
					return fmt.Errorf("failed to parse listing: %w", err)
				}
				point := listing.Location.ToPoint()
				locHash := point.Hash()
				queue[locHash] = true
			}
			return nil
		},
	)
	if err != nil {
		return err
	}

	log.Printf("Found %d unique geohashes to reverse-geocode", len(queue))
	//reverse-geocode each set and push it into a new bulk-insert list
	bulk := make([]domain.City, 0)
	for hash := range queue {
		lat, lng := geohash.Decode(hash)
		loc := domain.Location{Latitude: lat, Longitude: lng}
		res, err := geocode.ReverseGeocode(loc.ToPoint(), geoapify.ResultLevelCity, ctx)
		if err != nil {
			return err
		}
		log.Printf("Found city %s (%f, %f) %s\n", res.FormattedName, lat, lng, hash)
		resolvedCityLoc := domain.Location{
			Latitude:  res.Latitude,
			Longitude: res.Longitude,
		}
		bulk = append(bulk, domain.City{
			Location:    resolvedCityLoc,
			Name:        res.Name,
			SubDivision: res.State,
			Timezone:    res.Timezone.Name,
			FullName:    res.FormattedName,
			Hash:        hash,
		})
	}
	db.CreateInBatches(bulk, 100)
	return nil
}

func getFileHandle(
	path string,
) (*os.File, error) {
	if fs.ValidPath(path) {
		file, err := os.Open(path)
		if err != nil {
			return file, err
		}
		return file, nil
	}
	return nil, fmt.Errorf("path %s is not valid", path)
}

func rewindFileHandle(file *os.File) error {
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("rewind input file: %w", err)
	}
	return nil
}
