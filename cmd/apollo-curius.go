package main

import (
	"log"
	"os"
	"time"

	"github.com/straightupjac/curius-search/pkg/apollo-curius"
	"github.com/straightupjac/curius-search/pkg/apollo-curius/backend"
	"github.com/straightupjac/curius-search/pkg/curius"
)

func main() {
	// sources.ReadXMLFile()
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) != 0 && argsWithoutProg[0] == "--export" {
		curius.SaveCuriusArticles()
	} else {
		backend.InitializeFilesAndData()
		//we call ticker to refresh inverted index regularly once every 3 days however
		//for convenience we often want to do a refresh "on start" so we add this here too
		backend.RefreshInvertedIndex()
		log.Println("Refreshing inverted index on launch: ")
		// two days in miliseconds
		// once every three days, takes all the records, pulls from the data sources,
		ticker := time.NewTicker(3 * 24 * time.Hour)
		done := make(chan bool)
		go func() {
			for {
				select {
				case <-done:
					return
				case t := <-ticker.C:
					log.Println("Refreshing inverted index at: ", t)
					backend.RefreshInvertedIndex()
				}
			}
		}()
		//server and the pipeline should run on concurrent threads, called regularly, for now manually do it
		//start the server on a concurrent thread so that when we need to refresh the inverted index, this happens on
		//different threads
		// backend.RefreshInvertedIndex()
		apollo.Start()
	}

}
