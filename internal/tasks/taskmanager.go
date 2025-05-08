package tasks

import (
	"encoding/json"
	"gartifacts/internal/itemplanner"
	"github.com/thestuckster/gopherfacts/pkg/clients"
	"github.com/thestuckster/gopherfacts/pkg/items"
	"log"
	"os"
)

// TODO: need to come up with a better task schema.
// TODO: need to look up map coordinates for each associated task / resource
// TODO: need to add pulling / depositing into bank tasks, see first todo
// TODO: honestly, i might need some DI. this is getting dumb.
// TODO: git gudder
type Task struct {
	Action string
	X      int
	Y      int
}

func NewTask(node itemplanner.ItemPlanNode) Task {
	t := Task{}

	if node.Meta.Craft.Skill == "" {
		t.Action = node.Meta.Subtype
	} else {
		t.Action = node.Meta.Craft.Skill
	}
	return t
}

// TODO: check if you CAN actually gather / craft the thing
func BuildTaskQueue(characterName, desiredItem string,
	itemsByName, itemsByCode map[string]items.ItemMetaData,
	client *clients.GopherFactClient) []Task {

	log.Printf("Building a task queue for %s\n", desiredItem)
	tasks := make([]Task, 0)

	desiredItemNode, err := itemplanner.BuildItemPlanGraph(desiredItem, itemsByName, itemsByCode)
	if err != nil {
		log.Fatal(err)
	}

	craftingOrder := desiredItemNode.BuildCraftingOrder()
	for _, resource := range craftingOrder {
		log.Printf("Building task for %s\n", resource.Name)

		inventoryQuantity := checkInventoryQuantity(characterName, resource.Code, client)
		bankQuantity := checkBankQuantity(characterName, resource.Code, client)

		if inventoryQuantity == 0 && bankQuantity == 0 {
			log.Printf("No existing stock of %s in inventory or bank", resource.Name)
			for range resource.Amount {
				tasks = append(tasks, NewTask(*resource))
			}
			continue
		}

		if inventoryQuantity > 0 {
			if inventoryQuantity > resource.Amount {
				log.Println("We already have more than enough of this item in our inventory")
				continue
			}

			neededAmount := resource.Amount - inventoryQuantity
			for range neededAmount {
				tasks = append(tasks, NewTask(*resource))
			}
		}

		if bankQuantity > 0 {
			if bankQuantity > resource.Amount {
				log.Println("We already have more than enough of this item in the bank")
				//TODO: implement bank actions for the queue. I'm not dealing with that tonight
				continue
			}

			neededAmount := resource.Amount - bankQuantity
			for range neededAmount {
				tasks = append(tasks, NewTask(*resource))
			}
		}
	}

	return tasks
}

func characterCanFarmResource(characterName, skill string, requiredLevel int, client *clients.GopherFactClient) bool {
	char, err := client.CharacterClient.GetCharacterInfo(characterName)
	if err != nil {
		log.Fatal(err)
	}

	switch skill {
	case "mining":
		return char.MiningLevel >= requiredLevel
	case "weaponcrafting":
		return char.WeaponcraftingLevel >= requiredLevel
	case "woodcutting":
		return char.WoodcuttingLevel >= requiredLevel
	case "fishing":
		return char.FishingLevel >= requiredLevel
	case "gearcrafting":
		return char.GearcraftingLevel >= requiredLevel
	case "jewelrycrafting":
		return char.JewelrycraftingLevel >= requiredLevel
	case "cooking":
		return char.CookingLevel >= requiredLevel
	default:
		return false
	}
}

// boy what a name that one is
func checkInventoryQuantity(characterName, itemCode string, client *clients.GopherFactClient) int {

	char, err := client.CharacterClient.GetCharacterInfo(characterName)
	if err != nil {
		log.Fatal(err)
	}

	for _, inventoryItem := range char.Inventory {
		if itemCode == inventoryItem.Code {
			return inventoryItem.Quantity
		}
	}

	return 0
}

func checkBankQuantity(characterName, itemCode string, client *clients.GopherFactClient) int {
	bankItems, err := client.AccountClient.GetAllBankItems()
	if err != nil {
		log.Fatal(err)
	}

	for _, bankItem := range bankItems {
		if itemCode == bankItem.Code {
			return bankItem.Quantity
		}
	}

	return 0
}

// debugGraph prints the current crafting order to a json file for visual debugging and inspection
func debugGraph(items []*itemplanner.ItemPlanNode) {
	log.Println("Debugging graph by writing crafting order to a json file")
	jsonData, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("itemGraph.json", jsonData, 0644)
	if err != nil {
		panic(err)
	}
}
