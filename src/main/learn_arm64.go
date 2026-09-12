package main

/*
#include "arm64.h"
*/
import "C"

import (
	"fmt"
	"time"
)

const (
	dimension  = 128
	iterations = 10_000_000
)

// --------------------------------------------------
// Pure Go int32 implementation
// --------------------------------------------------

func vectorDistanceGoInt32(
	vec1 *[dimension]C.int32_t,
	vec2 *[dimension]C.int32_t,
) int64 {

	var distance int64

	for i := 0; i < dimension; i++ {
		diff := int64(vec1[i]) - int64(vec2[i])
		distance += diff * diff
	}

	return distance
}

// --------------------------------------------------
// Pure Go int8 implementation
// --------------------------------------------------

func vectorDistanceGoInt8(
	vec1 *[dimension]C.int8_t,
	vec2 *[dimension]C.int8_t,
) int64 {

	var distance int64

	for i := 0; i < dimension; i++ {
		diff := int64(vec1[i]) - int64(vec2[i])
		distance += diff * diff
	}

	return distance
}

func main() {

	// ==================================================
	// Create vectors
	// ==================================================

	vec1 := [dimension]C.int32_t{}
	vec2 := [dimension]C.int32_t{}

	vec1_8 := [dimension]C.int8_t{}
	vec2_8 := [dimension]C.int8_t{}

	for i := 0; i < dimension; i++ {

		// int32 vectors
		vec1[i] = C.int32_t(i)
		vec2[i] = C.int32_t(2 * i)

		// int8 vectors
		vec1_8[i] = C.int8_t(i)
		vec2_8[i] = C.int8_t(2 * i)
	}

	// ==================================================
	// Correctness
	// ==================================================

	goInt32 := vectorDistanceGoInt32(&vec1, &vec2)

	asmInt32 := C.vector_distance(
		&vec1[0],
		&vec2[0],
	)

	goInt8 := vectorDistanceGoInt8(&vec1_8, &vec2_8)

	asmInt8 := C.vector_distance_8bit(
		&vec1_8[0],
		&vec2_8[0],
	)

	fmt.Println("Correctness")
	fmt.Println("--------------------------------")

	fmt.Println("INT32")
	fmt.Println("Go:  ", goInt32)
	fmt.Println("ASM: ", asmInt32)

	if goInt32 != int64(asmInt32) {
		panic("INT32 results do not match")
	}

	fmt.Println("✓ INT32 results match!")
	fmt.Println()

	fmt.Println("INT8")
	fmt.Println("Go:  ", goInt8)
	fmt.Println("ASM: ", asmInt8)

	if goInt8 != int64(asmInt8) {
		panic("INT8 results do not match")
	}

	fmt.Println("✓ INT8 results match!")

	// ==================================================
	// Benchmark INT32 Go
	// ==================================================

	start := time.Now()

	var goInt32Sum int64

	for i := 0; i < iterations; i++ {
		goInt32Sum += vectorDistanceGoInt32(
			&vec1,
			&vec2,
		)
	}

	goInt32Time := time.Since(start)

	// ==================================================
	// Benchmark INT32 ARM64
	// ==================================================

	start = time.Now()

	var asmInt32Sum int64

	for i := 0; i < iterations; i++ {
		asmInt32Sum += int64(
			C.vector_distance(
				&vec1[0],
				&vec2[0],
			),
		)
	}

	asmInt32Time := time.Since(start)

	// ==================================================
	// Benchmark INT8 Go
	// ==================================================

	start = time.Now()

	var goInt8Sum int64

	for i := 0; i < iterations; i++ {
		goInt8Sum += vectorDistanceGoInt8(
			&vec1_8,
			&vec2_8,
		)
	}

	goInt8Time := time.Since(start)

	// ==================================================
	// Benchmark INT8 ARM64
	// ==================================================

	start = time.Now()

	var asmInt8Sum int64

	for i := 0; i < iterations; i++ {
		asmInt8Sum += int64(
			C.vector_distance_8bit(
				&vec1_8[0],
				&vec2_8[0],
			),
		)
	}

	asmInt8Time := time.Since(start)

	// ==================================================
	// Checksums
	// ==================================================

	fmt.Println()
	fmt.Println("Checksums")
	fmt.Println("--------------------------------")

	fmt.Println("Go INT32:  ", goInt32Sum)
	fmt.Println("ASM INT32: ", asmInt32Sum)
	fmt.Println("Go INT8:   ", goInt8Sum)
	fmt.Println("ASM INT8:  ", asmInt8Sum)

	// ==================================================
	// INT32 Results
	// ==================================================

	fmt.Println()
	fmt.Println("================================")
	fmt.Println("INT32")
	fmt.Println("================================")

	fmt.Println()
	fmt.Println("Go:")
	fmt.Println("  Total:    ", goInt32Time)
	fmt.Println(
		"  Per call: ",
		goInt32Time/time.Duration(iterations),
	)

	fmt.Println()
	fmt.Println("ARM64:")
	fmt.Println("  Total:    ", asmInt32Time)
	fmt.Println(
		"  Per call: ",
		asmInt32Time/time.Duration(iterations),
	)

	fmt.Printf(
		"\nARM64 speedup vs Go: %.2fx\n",
		float64(goInt32Time)/float64(asmInt32Time),
	)

	// ==================================================
	// INT8 Results
	// ==================================================

	fmt.Println()
	fmt.Println("================================")
	fmt.Println("INT8")
	fmt.Println("================================")

	fmt.Println()
	fmt.Println("Go:")
	fmt.Println("  Total:    ", goInt8Time)
	fmt.Println(
		"  Per call: ",
		goInt8Time/time.Duration(iterations),
	)

	fmt.Println()
	fmt.Println("ARM64:")
	fmt.Println("  Total:    ", asmInt8Time)
	fmt.Println(
		"  Per call: ",
		asmInt8Time/time.Duration(iterations),
	)

	fmt.Printf(
		"\nARM64 speedup vs Go: %.2fx\n",
		float64(goInt8Time)/float64(asmInt8Time),
	)

	// ==================================================
	// INT8 vs INT32
	// ==================================================

	fmt.Println()
	fmt.Println("================================")
	fmt.Println("INT8 vs INT32")
	fmt.Println("================================")

	fmt.Printf(
		"\nINT8 ARM64 vs INT32 ARM64: %.2fx\n",
		float64(asmInt32Time)/float64(asmInt8Time),
	)
}
