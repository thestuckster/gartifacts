package main

import (
	"gartifacts/internal"
	"github.com/thestuckster/gopherfacts/pkg/clients"
	"github.com/thestuckster/gopherfacts/pkg/items"
	"log"
)

const artifactsToken = "ARTIFACTS_TOKEN"

func main() {

	mainCharacter := "Main"

	allItemsByName, _, err := fetchAllItemInformation()
	if err != nil {
		log.Fatal(err)
	}

	//debug
	woodStaff := allItemsByName["Wooden Staff"]
	log.Println(woodStaff)

	///ensure env vars are set
	internal.LoadEnvironment()
	apiToken := internal.GetEnvVar(artifactsToken)

	client := clients.NewClient(&apiToken)
	_, err = client.EasyClient.MoveToChickens(mainCharacter)
	if err != nil {
		log.Fatal(err)
	}

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
