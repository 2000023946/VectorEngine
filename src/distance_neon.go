package src

// SquaredDistanceNEON computes squared distance on quantized int32 vectors,
// accumulating in int64 to avoid overflow on high-dimensional data.
func SquaredDistanceNEON(query []int8, vector []int8) int64
