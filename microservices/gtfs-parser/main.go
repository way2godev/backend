package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"way2go/bootstrap"
	"way2go/domain/entities"
	"way2go/infraestructure/database"
	"way2go/microservices/gtfs-parser/constants"
	"way2go/microservices/gtfs-parser/parsers"

	"github.com/mholt/archiver"
)



func main() {
	bootstrap.Init()

	db := database.GetDB()

	var datasources []entities.Datasource
	db.Find(&datasources)
	log.Printf("Found in total %d datasources\n", len(datasources))

	// Create the data directory if it doesn't exist
	if _, err := os.Stat(constants.WORKDIR); os.IsNotExist(err) {
		os.Mkdir(constants.WORKDIR, 0777)
		log.Println("Data directory created")
	} else {
		log.Println("Data directory already exists")
		removeItemsFromWorkdir()
	}

	for _, source := range datasources {
		log.Printf("Processing source %s\n", source.Name)
		log.Printf("Source data: %v\n", source)

		// Download the zip file
		err := downloadZip(source.Url)
		if err != nil {
			log.Printf("Error downloading zip file: %v\n", err)
			log.Printf("Skipping source %s\n", source.Name)
			continue
		}

		unzipFile("file.zip")

		parsers.Agencies()
		parsers.Routes()
		parsers.Stops()

		removeItemsFromWorkdir()
		log.Printf("Finished processing source %s\n", source.Name)
	}

}

func removeItemsFromWorkdir() {
	items, err := os.ReadDir(constants.WORKDIR)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}
	for _, item := range items {
		err := os.RemoveAll(fmt.Sprintf("%s/%s", constants.WORKDIR, item.Name()))
		if err != nil {
			log.Fatalf("Error removing item: %v", err)
		}
	}
}

func downloadZip(url string) error {
	log.Printf("Downloading file from %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error downloading file: %v", err)
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(fmt.Sprintf("%s/%s", constants.WORKDIR, "file.zip"))
	if err != nil {
		log.Printf("Error creating file: %v", err)
		return err
	}
	defer out.Close()

	// Write the content
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Printf("Error writing file: %v", err)
		return err
	}
	log.Printf("Downloaded file from %s\n", url)
	return nil
}

func unzipFile(filename string) {
	log.Printf("Unzipping file %s\n", filename)

	// Unzip the file
	err := archiver.Unarchive(constants.WORKDIR+"/"+filename, constants.WORKDIR)
	if err != nil {
		log.Fatalf("Error unzipping file: %v", err)
	}

	log.Printf("Unzipped file %s\n", filename)
}
