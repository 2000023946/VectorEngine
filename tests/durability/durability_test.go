package src

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
	"vectorengine/src"
)

const (
	durabilityDimension = 128
)

func generateTestVector() []float64 {
	vector := make([]float64, durabilityDimension)

	for i := range vector {
		vector[i] = rand.Float64()
	}

	return vector
}

func TestDurability(t *testing.T) {

	sizes := []int{
		10_000,
		100_000,
		1_000_000,
	}

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("       VECTOR ENGINE DURABILITY TEST")
	fmt.Println("========================================")
	fmt.Println()

	for _, size := range sizes {

		t.Run(fmt.Sprintf("%d_vectors", size), func(t *testing.T) {

			fmt.Println()
			fmt.Println("----------------------------------------")
			fmt.Printf("Testing %d vectors\n", size)
			fmt.Println("----------------------------------------")

			// Start clean.
			os.Remove("vectorengine.dat")

			// --------------------------------------------------
			// Initialize API
			// --------------------------------------------------

			fmt.Println("[1/5] Initializing API...")

			api, err := src.NewAPI()
			if err != nil {
				t.Fatalf("failed to initialize API: %v", err)
			}

			fmt.Println("      API initialized.")

			// --------------------------------------------------
			// Insert
			// --------------------------------------------------

			fmt.Printf("[2/5] Inserting %d vectors...\n", size)

			insertStart := time.Now()

			nextProgress := 10

			for i := 0; i < size; i++ {

				vector := generateTestVector()

				if err := api.Insert(i, vector); err != nil {
					t.Fatalf(
						"insert failed at vector %d: %v",
						i,
						err,
					)
				}

				// Print progress every 10%.
				progress := ((i + 1) * 100) / size

				if progress >= nextProgress {
					fmt.Printf(
						"      Insert progress: %d%% (%d/%d)\n",
						nextProgress,
						i+1,
						size,
					)

					nextProgress += 10
				}
			}

			insertDuration := time.Since(insertStart)

			fmt.Printf(
				"      Insert complete: %v\n",
				insertDuration,
			)

			// --------------------------------------------------
			// Persistence
			// --------------------------------------------------

			fmt.Println("[3/5] Waiting for asynchronous persistence...")

			persistenceStart := time.Now()

			api.WaitForPersistence()

			persistenceDuration := time.Since(persistenceStart)

			fmt.Printf(
				"      Persistence complete: %v\n",
				persistenceDuration,
			)

			// --------------------------------------------------
			// Simulate crash
			// --------------------------------------------------

			fmt.Println("[4/5] Simulating shutdown/crash...")

			api = nil

			fmt.Println("      Engine discarded.")

			// --------------------------------------------------
			// Reboot
			// --------------------------------------------------

			fmt.Println("[5/5] Rebooting from disk...")

			rebootStart := time.Now()

			recoveredAPI, err := src.NewAPI()

			rebootDuration := time.Since(rebootStart)

			if err != nil {
				t.Fatalf(
					"reboot failed: %v",
					err,
				)
			}

			fmt.Printf(
				"      Reboot complete: %v\n",
				rebootDuration,
			)

			// --------------------------------------------------
			// Verify durability
			// --------------------------------------------------

			recoveredCount := recoveredAPI.Count()

			if recoveredCount != size {
				fmt.Printf(
					"      ❌ DURABILITY FAILED\n",
				)

				t.Fatalf(
					"expected %d vectors, recovered %d",
					size,
					recoveredCount,
				)
			}

			fmt.Printf(
				"      ✅ DURABILITY PASSED\n",
			)

			fmt.Println()
			fmt.Println("Results:")
			fmt.Printf(
				"  Vectors:             %d\n",
				size,
			)
			fmt.Printf(
				"  Insert time:         %v\n",
				insertDuration,
			)
			fmt.Printf(
				"  Persistence wait:    %v\n",
				persistenceDuration,
			)
			fmt.Printf(
				"  Reboot time:         %v\n",
				rebootDuration,
			)
			fmt.Printf(
				"  Recovered vectors:   %d/%d\n",
				recoveredCount,
				size,
			)

			// Clean up.
			if err := recoveredAPI.Reset(); err != nil {
				t.Fatalf(
					"failed to reset after durability test: %v",
					err,
				)
			}
		})
	}

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("       DURABILITY TEST COMPLETE")
	fmt.Println("========================================")
	fmt.Println()
}
