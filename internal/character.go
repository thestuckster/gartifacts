package internal

import (
	"fmt"
	"github.com/thestuckster/gopherfacts/pkg/clients"
	"log"
	"sync"
	"time"
)

func ChickenLoop(characterName string, client *clients.GopherFactClient, stopChannel chan bool, wg *sync.WaitGroup) {

	defer wg.Done()

	characterDetails, err := client.CharacterClient.GetCharacterInfo(characterName)
	if err != nil {
		log.Fatalf("Error getting character details: %v\n", err)
	}

	log.Println("Chicken loop started")
	for {
		// None blocking channel listening
		select {
		case stop, ok := <-stopChannel:
			// Handle the received value
			fmt.Printf("Received: %v, Channel open: %v\n", stop, ok)
			if stop {
				log.Println("Chicken Loop manually stopped")
				return
			}
		default:
			// This executes immediately if stopChannel doesn't have a value ready
			//fmt.Println("No value available, not blocking")
		}

		//check inventory, if full, move to bank and deposit
		checkAndDumpInventory(characterName, *characterDetails, client)

		//move to the chickens
		log.Println("Moving to chickens")
		_, err := client.EasyClient.MoveToChickens(characterName)
		if err != nil {
			panic(err)
		}

		//fight
		fightResponse, err := client.CharacterClient.Fight(characterName)
		if err != nil {
			log.Fatalf("Error fighting %v: %v\n", characterName, err)
		}
		log.Printf("Chicken %v fight %v\n", characterName, fightResponse)
		time.Sleep(time.Duration(fightResponse.Cooldown.RemainingSeconds) * time.Second)

		//update characterDetails
		characterDetails = &fightResponse.Character

		//rest if healthThreshold is below
		shouldTakeRest := shouldRest(*characterDetails)
		for shouldTakeRest == true {
			log.Printf("Character health is at %d / %d... resting", characterDetails.HP, characterDetails.MaxHP)
			restResponse, err := client.EasyClient.Rest(characterName)
			if err != nil {
				log.Fatalf("Error resting %v: %v\n", characterName, err)
			}

			characterDetails.HP += restResponse.HpRestored
			log.Printf("Gained %d HP from rest. Health now at %d\n", restResponse.HpRestored, characterDetails.HP)
			shouldTakeRest = shouldRest(*characterDetails)
			log.Printf("Should we take another rest? %v\n", shouldTakeRest)
		}
	}
}

func shouldRest(characterDetails clients.CharacterSchema) bool {
	return characterDetails.HP < characterDetails.MaxHP
}

func checkAndDumpInventory(characterName string, characterDetails clients.CharacterSchema, client *clients.GopherFactClient) {
	for _, slot := range characterDetails.Inventory {
		if isStackFull(slot) {
			_, err := client.EasyClient.MoveToBank(characterName)
			if err != nil {
				log.Fatalf("Error moving %v: %v\n", characterName, err)
			}

			log.Printf("Moved to bank... Storing %d of %s\n", slot.Quantity, slot.Code)
			_, err = client.EasyClient.DepositIntoBank(characterName, slot.Code, slot.Quantity)
			if err != nil {
				log.Fatalf("Error moving %v: %v\n", characterName, err)
			}
		}
	}
}

func isStackFull(slot clients.InventorySlot) bool {
	return slot.Quantity >= 100
}
