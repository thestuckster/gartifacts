package itemplanner

import (
	"github.com/thestuckster/gopherfacts/pkg/items"
	"testing"
)

func TestBuildItemPlanGraph(t *testing.T) {

	itemsByName, itemsByCode := buildSimpleItemMaps()
	g, err := BuildItemPlanGraph("Wooden Staff", itemsByName, itemsByCode)
	if err != nil {
		t.Fatal(err)
	}

	if len(g.RequiredComponents) < 2 {
		t.Fatalf("Required components count is incorrect. Expected: 2, Actual: %d", len(g.RequiredComponents))
	}
}

func TestFindItemInGraph_Name(t *testing.T) {
	itemsByName, itemsByCode := buildSimpleItemMaps()
	g, err := BuildItemPlanGraph("Wooden Staff", itemsByName, itemsByCode)
	if err != nil {
		t.Fatal(err)
	}

	ashWood := g.FindItemInGraph("Ash Wood", "")
	if ashWood == nil {
		t.Fatalf("Ash Wood not found in graph")
	}
}

func TestFindItemInGraph_Code(t *testing.T) {
	itemsByName, itemsByCode := buildSimpleItemMaps()
	g, err := BuildItemPlanGraph("Wooden Staff", itemsByName, itemsByCode)
	if err != nil {
		t.Fatal(err)
	}

	ashWood := g.FindItemInGraph("", "ash_wood")
	if ashWood == nil {
		t.Fatalf("Ash Wood not found in graph")
	}
}

func TestBuildItemCraftingOrder(t *testing.T) {
	itemsByName, itemsByCode := buildSimpleItemMaps()
	g, err := BuildItemPlanGraph("Wooden Staff", itemsByName, itemsByCode)
	if err != nil {
		t.Fatal(err)
	}

	craftingOrder := g.BuildCraftingOrder()
	if len(craftingOrder) < 3 {
		t.Fatalf("Crafting order count is incorrect. Expected: 3, Actual: %d", len(craftingOrder))
	}
}

func buildSimpleItemMaps() (map[string]items.ItemMetaData, map[string]items.ItemMetaData) {

	woodenStickCraft := items.CraftItem{
		Code:     "wooden_stick",
		Quantity: 1,
	}

	ashWoodCraft := items.CraftItem{
		Code:     "ash_wood",
		Quantity: 4,
	}

	//not sure why it didn't like the [...]items.CraftItem {woodenStickCraft, ashWoodCraft} method
	requiredItems := make([]items.CraftItem, 0)
	requiredItems = append(requiredItems, woodenStickCraft)
	requiredItems = append(requiredItems, ashWoodCraft)

	woodenStaff := items.ItemMetaData{
		Name:  "Wooden Staff",
		Code:  "wooden_staff",
		Level: 1,
		Type:  "Weapon",
		Craft: items.Craft{
			Skill: "Weaponcrafting",
			Items: requiredItems,
		},
	}

	woodenStick := items.ItemMetaData{
		Name:  "Wooden Stick",
		Code:  "wooden_stick",
		Level: 1,
		Type:  "Weapon",
	}

	ashWood := items.ItemMetaData{
		Name:  "Ash Wood",
		Code:  "ash_wood",
		Level: 1,
		Type:  "Resource",
	}

	itemsByName := make(map[string]items.ItemMetaData)
	itemsByName[woodenStick.Name] = woodenStick
	itemsByName[ashWood.Name] = ashWood
	itemsByName[woodenStaff.Name] = woodenStaff

	itemsByCode := make(map[string]items.ItemMetaData)
	itemsByCode[woodenStick.Code] = woodenStick
	itemsByCode[ashWood.Code] = ashWood
	itemsByCode[woodenStaff.Code] = woodenStaff

	return itemsByName, itemsByCode
}
