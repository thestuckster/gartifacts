package itemplanner

import (
	"errors"

	"github.com/thestuckster/gopherfacts/pkg/items"
)

type ItemPlanNode struct {
	Meta   items.ItemMetaData
	Name   string
	Code   string
	Amount int

	RequiredComponents []*ItemPlanNode
}

// BuildItemPlanGraph Creates a graph representation of an item and all the items required to craft it
func BuildItemPlanGraph(desiredItemName string, itemsByName, itemsByCode map[string]items.ItemMetaData) (*ItemPlanNode, error) {
	desiredItem, ok := itemsByName[desiredItemName]
	if !ok {
		return nil, errors.New(desiredItemName + " item does not exist")
	}

	root := ItemPlanNode{
		Meta:               desiredItem,
		Name:               desiredItem.Name,
		Code:               desiredItem.Code,
		Amount:             1,
		RequiredComponents: []*ItemPlanNode{},
	}

	for _, itemToCraft := range desiredItem.Craft.Items {
		name := itemsByCode[itemToCraft.Code].Name
		next, err := BuildItemPlanGraph(name, itemsByName, itemsByCode)
		if err != nil {
			return nil, err
		}

		next.Amount = itemToCraft.Quantity
		root.RequiredComponents = append(root.RequiredComponents, next)
	}

	return &root, nil
}

// FindItemInGraph performs a depth-first search to find a specific ItemPlanNode by name or code
// Returns the found node or nil if not found
func (n *ItemPlanNode) FindItemInGraph(targetName, targetCode string) *ItemPlanNode {
	if n == nil {
		return nil
	}

	// Check if current node matches the target
	if n.Name == targetName || n.Code == targetCode {
		return n
	}

	// Recursively search through children
	for _, child := range n.RequiredComponents {
		found := child.FindItemInGraph(targetName, targetCode)
		if found != nil {
			return found
		}
	}

	return nil
}

// BuildCraftingOrder performs a breadth-first search traversal of the ItemPlanGraph
// It returns a slice of ItemPlanNode pointers in the order they were visited
func (n *ItemPlanNode) BuildCraftingOrder() []*ItemPlanNode {
	if n == nil {
		return []*ItemPlanNode{}
	}

	// Initialize queue with root node
	queue := []*ItemPlanNode{n}
	visited := make(map[*ItemPlanNode]bool)
	result := []*ItemPlanNode{n}

	// Process nodes in queue
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// Add all unvisited children to queue
		for _, child := range current.RequiredComponents {
			if !visited[child] {
				visited[child] = true
				queue = append(queue, child)
				result = append(result, child)
			}
		}
	}

	return result
}
