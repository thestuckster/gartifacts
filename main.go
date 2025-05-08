package main

import (
	"gartifacts/internal"
	"gartifacts/internal/tasks"
	"github.com/thestuckster/gopherfacts/pkg/clients"
	"github.com/thestuckster/gopherfacts/pkg/items"
	"github.com/thestuckster/gopherfacts/pkg/maps"
	"log"
	"sync"
)

const artifactsToken = "ARTIFACTS_TOKEN"

func main() {
	log.Println("Starting gartifacts")
	mainCharacter := "Main"

	allItemsByName, allItemsByCode, err := fetchAllItemInformation()
	if err != nil {
		log.Fatal(err)
	}

	//tilesByResourceCode, err := fetchAllMapInformation()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println(len(tilesByResourceCode))
	//
	////debug
	//woodStaff := allItemsByName["Wooden Staff"]
	//log.Println(woodStaff)
	//

	///ensure env vars are set
	internal.LoadEnvironment()
	apiToken := internal.GetEnvVar(artifactsToken)

	client := clients.NewClient(&apiToken)

	//boolChan := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(2)

	log.Println("Starting chicken loop")
	//go internal.ChickenLoop(mainCharacter, client, boolChan, &wg)
	tasks := tasks.BuildTaskQueue(mainCharacter, "Copper Dagger", allItemsByName, allItemsByCode, client)
	log.Println(tasks)
	wg.Wait()
}

func fetchAllItemInformation() (map[string]items.ItemMetaData, map[string]items.ItemMetaData, error) {
	log.Println("Fetching all item information")
	allItems, err := items.GetAllItemData()
	if err != nil {
		return nil, nil, err
	}

	itemsByName := make(map[string]items.ItemMetaData)
	itemsByCode := make(map[string]items.ItemMetaData)

	for _, item := range allItems {
		itemsByName[item.Name] = item
		itemsByCode[item.Code] = item
	}

	return itemsByName, itemsByCode, nil
}

func fetchAllMapInformation() (tilesByResourceCode map[string][]maps.MapData, err error) {
	log.Println("Fetching all map information")
	mapTiles, err := maps.GetAllMapData()
	if err != nil {
		return nil, err
	}

	tilesByResourceCode = make(map[string][]maps.MapData)
	for _, tile := range mapTiles {
		resource := tile.Content.Code
		if _, ok := tilesByResourceCode[resource]; !ok {
			tilesByResourceCode[resource] = []maps.MapData{tile}
		} else {
			tilesByResourceCode[resource] = append(tilesByResourceCode[resource], tile)
		}
	}

	return tilesByResourceCode, err
}
