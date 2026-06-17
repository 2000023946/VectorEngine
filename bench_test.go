package vectorengine

import "testing"

func BenchmarkInsert(b *testing.B) {
	vec := make([]float32, 128)
	for i := range vec {
		vec[i] = 1.0
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		store := NewVectorStore(128, 1) // tiny capacity every iteration

		for j := 0; j < 5000; j++ {
			store.insert(vec)
		}
	}
}
func BenchmarkSearch(b *testing.B) {

	store := NewVectorStore(128, 100000)

	vec := make([]float32, 128)
	for i := range vec {
		vec[i] = float32(i)
	}

	// preload data
	for i := 0; i < 100000; i++ {
		store.insert(vec)
	}

	query := make([]float32, 128)
	for i := range query {
		query[i] = 0.5
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = store.Search(query)
	}
}
