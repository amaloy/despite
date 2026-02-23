package core

import (
	"testing"
)

func TestToDSChar(t *testing.T) {
	tests := []struct {
		input    int
		expected rune
	}{
		{0, 32},   // space
		{1, 33},   // !
		{65, 97},  // a
		{90, 122}, // z
	}

	for _, test := range tests {
		result := toDSChar(test.input)
		if result != test.expected {
			t.Errorf("toDSChar(%d) = %d; want %d", test.input, result, test.expected)
		}
	}
}

func TestBuildMainMap(t *testing.T) {
	m, err := BuildMainMap()
	if err != nil {
		t.Fatalf("BuildMainMap failed: %v", err)
	}
	if m == nil {
		t.Fatal("BuildMainMap returned nil")
	}
	if m.name != "lev01" {
		t.Errorf("Expected name 'lev01', got %s", m.name)
	}
	if m.width != standardMapWidth {
		t.Errorf("Expected width %d, got %d", standardMapWidth, m.width)
	}
	if m.height != standardMapHeight {
		t.Errorf("Expected height %d, got %d", standardMapHeight, m.height)
	}
	if m.xstart != 26 {
		t.Errorf("Expected xstart 26, got %d", m.xstart)
	}
	if m.ystart != 41 {
		t.Errorf("Expected ystart 41, got %d", m.ystart)
	}
	if len(m.tiles) != m.width {
		t.Errorf("Expected tiles length %d, got %d", m.width, len(m.tiles))
	}
	for x := range m.tiles {
		if len(m.tiles[x]) != m.height {
			t.Errorf("Expected row %d length %d, got %d", x, m.height, len(m.tiles[x]))
		}
		for y := range m.tiles[x] {
			if m.tiles[x][y] == nil {
				t.Errorf("Tile at (%d,%d) is nil", x, y)
			}
		}
	}
	if m.Players == nil {
		t.Error("players map is nil")
	}
}
