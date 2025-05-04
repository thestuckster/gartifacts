package maps

import (
	"testing"

	"github.com/thestuckster/gopherfacts/pkg/maps"
)

// auto generated test case
func TestFindClosetMapForResource(t *testing.T) {
	tests := []struct {
		name         string
		charX, charY int
		resourceCode string
		tiles        map[string][]maps.MapData
		want         *maps.MapData
		wantErr      bool
	}{
		{
			name:         "finds closest tile",
			charX:        0,
			charY:        0,
			resourceCode: "wood",
			tiles: map[string][]maps.MapData{
				"wood": {
					{X: 1, Y: 1},
					{X: 2, Y: 2},
					{X: 3, Y: 3},
				},
			},
			want:    &maps.MapData{X: 1, Y: 1},
			wantErr: false,
		},
		{
			name:         "returns error for missing resource",
			charX:        0,
			charY:        0,
			resourceCode: "stone",
			tiles: map[string][]maps.MapData{
				"wood": {
					{X: 1, Y: 1},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:         "returns error for empty tile list",
			charX:        0,
			charY:        0,
			resourceCode: "wood",
			tiles: map[string][]maps.MapData{
				"wood": {},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindClosetMapForResource(tt.charX, tt.charY, tt.resourceCode, tt.tiles)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindClosetMapForResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil && (got.X != tt.want.X || got.Y != tt.want.Y) {
				t.Errorf("FindClosetMapForResource() = %v, want %v", got, tt.want)
			}
		})
	}
}
