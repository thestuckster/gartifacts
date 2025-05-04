package maps

import (
	"errors"
	"github.com/thestuckster/gopherfacts/pkg/maps"
	"math"
)

func FindClosetMapForResource(charX, charY int, resourceCode string, tilesByResource map[string][]maps.MapData) (*maps.MapData, error) {

	if tiles, ok := tilesByResource[resourceCode]; !ok {
		return nil, errors.New("No map tiles(s) found for containing resource " + resourceCode)
	} else {

		if len(tiles) == 0 {
			return nil, errors.New("No map tiles(s) found for containing resource " + resourceCode)
		}

		minDistance := math.MaxInt64
		var closetTile *maps.MapData

		for _, tile := range tiles {
			d := manhattanDistance(charX, charY, tile.X, tile.Y)
			if d < minDistance {
				minDistance = d
				closetTile = &tile
			}
		}

		return closetTile, nil
	}
}

// manhattanDistance measures the distance between points along axes at right angles.
// d = | x2 - x1 | + | y2 - y1|
func manhattanDistance(x1, y1, x2, y2 int) int {
	distance := math.Abs(float64(x2)-float64(x1)) + math.Abs(float64(y2)-float64(y1))

	//this conversion always truncates towards 0 so 42.9 --> 42 instead of 43
	return int(distance)
}
