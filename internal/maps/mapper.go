package maps

import (
	"errors"
	"github.com/thestuckster/gopherfacts/pkg/maps"
)

func FindClosetMapForResource(charX, charY int, resourceCode string, tilesByResource map[string][]maps.MapData) (*maps.MapData, error) {

	if tiles, ok := tilesByResource[resourceCode]; !ok {
		return nil, errors.New("No map tiles(s) found for containing resource " + resourceCode)
	} else {

		if len(tiles) == 0 {
			return nil, errors.New("No map tiles(s) found for containing resource " + resourceCode)
		}

		//TODO: find closest point on x, y graph
		return &tiles[0], nil
	}
}
