package surrender

import (
	"image"
	"testing"
)

func TestParsePath(t *testing.T) {
	tests := []struct {
		name        string
		pathData    string
		expectedCmd []string
		expectedPts [][]image.Point
	}{
		{
			name:        "Valid_Path_1",
			pathData:    "M 100 200 L 200 100 L -100 -200",
			expectedCmd: []string{"M", "L", "L"},
			expectedPts: [][]image.Point{
				{{100, 200}},
				{{200, 100}},
				{{-100, -200}},
			},
		},
		{
			name:        "Valid_Path_2",
			pathData:    "M100 200L200 100L-100-200",
			expectedCmd: []string{"M", "L", "L"},
			expectedPts: [][]image.Point{
				{{100, 200}},
				{{200, 100}},
				{{-100, -200}},
			},
		},
		// Add more test cases here as needed
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path, err := ParsePath(tc.pathData)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(path.Commands) != len(tc.expectedCmd) {
				t.Fatalf("Expected %d commands, but got %d", len(tc.expectedCmd), len(path.Commands))
				return
			}

			for i, cmd := range path.Commands {
				if cmd.Type != tc.expectedCmd[i] {
					t.Errorf("Expected command type %s, but got %s", tc.expectedCmd[i], cmd.Type)
				}

				expectedPts := tc.expectedPts[i]
				if len(cmd.Points) != len(expectedPts) {
					t.Errorf("Expected %d points for command type %s, but got %d",
						len(expectedPts), cmd.Type, len(cmd.Points))
				}

				for j, pt := range cmd.Points {
					if j >= len(expectedPts) {
						t.Errorf("Unexpected extra point %v for command type %s", pt, cmd.Type)
						continue
					}

					expectedPt := expectedPts[j]
					if pt != expectedPt {
						t.Errorf("Expected point %v for command type %s, but got %v",
							expectedPt, cmd.Type, pt)
					}
				}
			}
		})
	}
}
