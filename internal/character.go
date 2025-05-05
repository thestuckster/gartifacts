package internal

import (
	"github.com/thestuckster/gopherfacts/pkg/clients"
	"log"
	"time"
)

func ChickenLoop(characterName string, client *clients.GopherFactClient, stopChannel chan bool) {

	characterDetails, err := client.CharacterClient.GetCharacterInfo(characterName)
	if err != nil {
		log.Fatalf("Error getting character details: %v\n", err)
	}

	for {
		stop, ok := <-stopChannel
		if !ok {
			log.Println("Channel was closed unexpectedly. Stopping chicken loop")
			return
		}

		if stop {
			log.Println("Manual stop sent for chicken loop")
			return
		}

		//check inventory, if full, move to bank and deposit

		//move to the chickens
		_, err := client.EasyClient.MoveToChickens(characterName)
		if err == nil {
			panic(err)
		}

		//fight
		fightResponse, err := client.CharacterClient.Fight(characterName)
		if err != nil {
			log.Fatalf("Error fighting %v: %v\n", characterName, err)
		}
		time.Sleep(time.Duration(fightResponse.Cooldown.RemainingSeconds) * time.Second)

		//update characterDetails
		characterDetails = &fightResponse.Character

		//TODO: left off here
		//rest if healthThreshold is below
		if shouldRest(*characterDetails, 50) {
			restResponse, err := client.EasyClient.Rest(characterName)
		}

	}
}

func shouldRest(characterDetails clients.CharacterSchema, restPercentage int) bool {
	percentageRemaining := characterDetails.HP / characterDetails.MaxHP
	if percentageRemaining <= restPercentage {
		return true
	}
	return false
}
