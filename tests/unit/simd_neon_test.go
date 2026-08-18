package tests

// import (
// 	"testing"
// 	"vectorengine/src"
// )

// func TestSquaredDistanceNEON(t *testing.T) {
// 	tests := []struct {
// 		name   string
// 		query  []int8
// 		vector []int8
// 		want   int64
// 	}{
// 		{
// 			name:   "basic distance",
// 			query:  []int8{1, 2, 3, 4, 5, 6, 7, 8},
// 			vector: []int8{4, 3, 2, 1, 8, 7, 6, 5},
// 			want:   40,
// 		},
// 		{
// 			name:   "identical vectors",
// 			query:  []int8{1, 2, 3, 4, 5, 6, 7, 8},
// 			vector: []int8{1, 2, 3, 4, 5, 6, 7, 8},
// 			want:   0,
// 		},
// 		{
// 			name:   "all differences positive",
// 			query:  []int8{1, 1, 1, 1, 1, 1, 1, 1},
// 			vector: []int8{3, 3, 3, 3, 3, 3, 3, 3},
// 			want:   32,
// 		},
// 		{
// 			name:   "all differences negative",
// 			query:  []int8{3, 3, 3, 3, 3, 3, 3, 3},
// 			vector: []int8{1, 1, 1, 1, 1, 1, 1, 1},
// 			want:   32,
// 		},
// 		{
// 			name:   "zeros",
// 			query:  []int8{0, 0, 0, 0, 0, 0, 0, 0},
// 			vector: []int8{0, 0, 0, 0, 0, 0, 0, 0},
// 			want:   0,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := src.SquaredDistanceNEON(tt.query, tt.vector)

// 			if got != tt.want {
// 				t.Fatalf(
// 					"expected %d, got %d",
// 					tt.want,
// 					got,
// 				)
// 			}
// 		})
// 	}
// }

// func TestSquaredDistanceNEON128(t *testing.T) {
// 	query := make([]int8, 128)
// 	vector := make([]int8, 128)

// 	vector[0] = 100
// 	vector[1] = 100

// 	got := src.SquaredDistanceNEON(query, vector)

// 	// 100² + 100² = 20,000
// 	const want int64 = 20_000

// 	if got != want {
// 		t.Fatalf("expected %d, got %d", want, got)
// 	}
// }
