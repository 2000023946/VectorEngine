package src

import (
	"encoding/binary"
	"fmt"
	"os"
	"runtime"
	"sort"
	"sync"
)

const (
	maxVectors = 2_000_000
	dimension  = 128

	// Number of independent search partitions.
	searchWorkers = 6

	// quantizeScale controls precision vs. overflow headroom.
	quantizeScale = 10000.0

	// Maximum number of pending persistence operations.
	//
	// If this queue becomes full, Insert() blocks until
	// the persistence worker makes space.
	persistenceQueueSize = 16

	// File used for persistence.
	persistenceFile = "vectorengine.dat"
)

type Result struct {
	ID       int
	Distance float64
}

type VectorEngine struct {
	// --------------------------------------------------
	// Flat contiguous QUANTIZED vector storage
	// --------------------------------------------------

	vectorData []int32

	// External ID for each vector.
	ids []int

	dimension int
	count     int

	// Protects:
	//
	// count
	// vectorData
	// ids
	// resetting
	//
	// RLock allows multiple concurrent searches.
	// Lock provides exclusive access for inserts/resets.
	mu sync.RWMutex

	// --------------------------------------------------
	// Persistence
	// --------------------------------------------------

	// Bounded queue of persistence requests.
	//
	// Inserts place save requests into this queue.
	// One persistence worker consumes the queue.
	persistenceQueue chan struct{}

	// Tracks persistence operations that have been
	// accepted by the system but have not finished.
	//
	// Reset waits for this to reach zero.
	saveWG sync.WaitGroup

	// True while Reset is in progress.
	//
	// Inserts received while resetting are rejected.
	resetting bool
}

func NewVectorEngine() *VectorEngine {
	ve := &VectorEngine{
		vectorData:       make([]int32, maxVectors*dimension),
		ids:              make([]int, maxVectors),
		dimension:        dimension,
		persistenceQueue: make(chan struct{}, 1000000),
	}

	go ve.persistenceWorker("vectorengine.dat")

	return ve
}

// --------------------------------------------------
// Persistence Worker
// --------------------------------------------------
//
// The worker waits for persistence requests.
//
// Every request causes the worker to take a consistent
// snapshot of RAM and write that snapshot to disk.
//
// Because there is only one worker, disk writes cannot
// race with each other.
//

func (ve *VectorEngine) persistenceWorker(path string) {
	for {
		// Wait for at least one persistence request.
		<-ve.persistenceQueue

		// We received one request.
		pending := 1

		// --------------------------------------------------
		// Drain everything currently waiting in the queue.
		// --------------------------------------------------
		//
		// Do not block waiting for more requests.
		// Take whatever is already available and make
		// one persistence snapshot for the entire batch.
		//
		for {
			select {
			case <-ve.persistenceQueue:
				pending++
			default:
				goto save
			}
		}

	save:
		// --------------------------------------------------
		// One save for the entire batch.
		// --------------------------------------------------

		ve.saveToDisk(path)

		// Every queued persistence request represented by
		// this batch is now complete.
		for i := 0; i < pending; i++ {
			ve.saveWG.Done()
		}
	}
}

// quantize converts a single float64 to its int32 quantized form.
func quantize(v float64) int32 {
	return int32(v * quantizeScale)
}

// quantizeVector quantizes a whole vector.
func quantizeVector(values []float64) []int32 {
	out := make([]int32, len(values))

	for i, v := range values {
		out[i] = quantize(v)
	}

	return out
}

func (ve *VectorEngine) Insert(id int, values []float64) error {
	if len(values) != ve.dimension {
		return fmt.Errorf(
			"dimension mismatch: expected %d, got %d",
			ve.dimension,
			len(values),
		)
	}

	// --------------------------------------------------
	// Protect insertion state.
	// --------------------------------------------------

	ve.mu.Lock()

	// Do not allow inserts during Reset.
	if ve.resetting {
		ve.mu.Unlock()

		return fmt.Errorf("reset in progress")
	}

	if ve.count >= maxVectors {
		ve.mu.Unlock()

		return fmt.Errorf(
			"vector capacity exceeded: maximum %d vectors",
			maxVectors,
		)
	}

	index := ve.count
	offset := index * ve.dimension

	// Quantize once at insert time.
	for i, v := range values {
		ve.vectorData[offset+i] = quantize(v)
	}

	ve.ids[index] = id
	ve.count++

	// --------------------------------------------------
	// Register persistence BEFORE releasing the lock.
	// --------------------------------------------------
	//
	// This save now belongs to the set of operations
	// that Reset must wait for.
	//

	ve.saveWG.Add(1)

	ve.mu.Unlock()

	// --------------------------------------------------
	// Enqueue persistence request.
	// --------------------------------------------------
	//
	// Normally this is extremely fast.
	//
	// If the queue is full, this blocks until the
	// persistence worker creates space.
	//

	ve.persistenceQueue <- struct{}{}

	return nil
}

// squaredDistance computes the distance between the quantized query
// and the vector stored at index using the NEON SIMD kernel.
func (ve *VectorEngine) squaredDistance(
	quantizedQuery []int32,
	index int,
) int64 {
	offset := index * ve.dimension

	stored := ve.vectorData[offset : offset+ve.dimension]

	return SquaredDistanceNEON(quantizedQuery, stored)
}

// searchRange searches one independent partition of the dataset.
func (ve *VectorEngine) searchRange(
	quantizedQuery []int32,
	start int,
	end int,
	k int,
) []Result {

	results := make([]Result, 0, k)

	for i := start; i < end; i++ {
		distance := ve.squaredDistance(quantizedQuery, i)

		result := Result{
			ID:       ve.ids[i],
			Distance: float64(distance),
		}

		if len(results) < k {
			results = append(results, result)
			continue
		}

		worst := 0

		for j := 1; j < k; j++ {
			if results[j].Distance > results[worst].Distance {
				worst = j
			}
		}

		if result.Distance < results[worst].Distance {
			results[worst] = result
		}
	}

	return results
}

func (ve *VectorEngine) Search(query []float64, k int) []Result {
	ve.mu.RLock()
	defer ve.mu.RUnlock()

	if len(query) != ve.dimension {
		return nil
	}

	if ve.count == 0 || k <= 0 {
		return []Result{}
	}

	if k > ve.count {
		k = ve.count
	}

	quantizedQuery := quantizeVector(query)

	workers := searchWorkers

	if workers > runtime.GOMAXPROCS(0) {
		workers = runtime.GOMAXPROCS(0)
	}

	if workers > ve.count {
		workers = ve.count
	}

	resultsPerWorker := make([][]Result, workers)

	chunkSize := (ve.count + workers - 1) / workers

	done := make(chan int, workers)

	for worker := 0; worker < workers; worker++ {
		start := worker * chunkSize
		end := start + chunkSize

		if end > ve.count {
			end = ve.count
		}

		if start >= ve.count {
			done <- worker
			continue
		}

		go func(worker, start, end int) {
			resultsPerWorker[worker] = ve.searchRange(
				quantizedQuery,
				start,
				end,
				k,
			)

			done <- worker
		}(worker, start, end)
	}

	for i := 0; i < workers; i++ {
		<-done
	}

	results := make([]Result, 0, k)

	for worker := 0; worker < workers; worker++ {
		for _, result := range resultsPerWorker[worker] {

			if len(results) < k {
				results = append(results, result)
				continue
			}

			worst := 0

			for j := 1; j < k; j++ {
				if results[j].Distance > results[worst].Distance {
					worst = j
				}
			}

			if result.Distance < results[worst].Distance {
				results[worst] = result
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Distance < results[j].Distance
	})

	return results
}

// --------------------------------------------------
// Reset
// --------------------------------------------------
//
// Reset:
//
// 1. Blocks new inserts.
// 2. Makes RAM logically empty.
// 3. Waits for all queued persistence operations.
// 4. Deletes the persisted database.
// 5. Allows inserts again.
//

func (ve *VectorEngine) Reset() error {

	// Establish reset boundary.
	ve.mu.Lock()

	if ve.resetting {
		ve.mu.Unlock()

		return fmt.Errorf("reset already in progress")
	}

	ve.resetting = true

	// RAM is immediately logically empty.
	ve.count = 0

	ve.mu.Unlock()

	// --------------------------------------------------
	// Wait until every persistence operation that was
	// accepted before Reset has finished.
	// --------------------------------------------------

	ve.saveWG.Wait()

	// --------------------------------------------------
	// All old persistence operations are finished.
	//
	// Now it is safe to remove the disk snapshot.
	// --------------------------------------------------

	err := os.Remove(persistenceFile)

	if err != nil && !os.IsNotExist(err) {

		ve.mu.Lock()
		ve.resetting = false
		ve.mu.Unlock()

		return fmt.Errorf(
			"failed to remove database snapshot: %w",
			err,
		)
	}

	// Reset is completely finished.
	ve.mu.Lock()
	ve.resetting = false
	ve.mu.Unlock()

	return nil
}

// --------------------------------------------------
// saveToDisk
// --------------------------------------------------
//
// Takes a consistent RAM snapshot and writes it to disk.
//
// The engine lock is held only while copying the data.
// Disk I/O happens after the lock is released.
//

func (ve *VectorEngine) saveToDisk(path string) {

	ve.mu.RLock()

	count := ve.count
	dimension := ve.dimension

	ids := make([]int, count)
	copy(ids, ve.ids[:count])

	vectorData := make([]int32, count*dimension)
	copy(
		vectorData,
		ve.vectorData[:count*dimension],
	)

	ve.mu.RUnlock()

	// --------------------------------------------------
	// RAM lock is now released.
	// Disk I/O does not block Insert/Search.
	// --------------------------------------------------

	file, err := os.Create(path)
	if err != nil {
		fmt.Printf(
			"saveToDisk: failed to create file: %v\n",
			err,
		)
		return
	}

	defer file.Close()

	if err := binary.Write(
		file,
		binary.LittleEndian,
		int64(count),
	); err != nil {
		fmt.Printf(
			"saveToDisk: failed to write count: %v\n",
			err,
		)
		return
	}

	if err := binary.Write(
		file,
		binary.LittleEndian,
		int64(dimension),
	); err != nil {
		fmt.Printf(
			"saveToDisk: failed to write dimension: %v\n",
			err,
		)
		return
	}

	for _, id := range ids {
		if err := binary.Write(
			file,
			binary.LittleEndian,
			int64(id),
		); err != nil {
			fmt.Printf(
				"saveToDisk: failed to write ID: %v\n",
				err,
			)
			return
		}
	}

	for _, value := range vectorData {
		if err := binary.Write(
			file,
			binary.LittleEndian,
			value,
		); err != nil {
			fmt.Printf(
				"saveToDisk: failed to write vector data: %v\n",
				err,
			)
			return
		}
	}

	if err := file.Sync(); err != nil {
		fmt.Printf(
			"saveToDisk: failed to sync file: %v\n",
			err,
		)
		return
	}
}

// --------------------------------------------------
// Reboot
// --------------------------------------------------
//
// Reboot is synchronous.
//
// The caller should invoke this before accepting
// requests from clients.
//

func (ve *VectorEngine) Reboot(path string) error {

	ve.mu.Lock()
	defer ve.mu.Unlock()

	file, err := os.Open(path)

	if err != nil {
		if os.IsNotExist(err) {
			ve.count = 0
			return nil
		}

		return fmt.Errorf(
			"failed to open database snapshot: %w",
			err,
		)
	}

	defer file.Close()

	var count int64
	var storedDimension int64

	if err := binary.Read(
		file,
		binary.LittleEndian,
		&count,
	); err != nil {
		return fmt.Errorf(
			"failed to read vector count: %w",
			err,
		)
	}

	if err := binary.Read(
		file,
		binary.LittleEndian,
		&storedDimension,
	); err != nil {
		return fmt.Errorf(
			"failed to read dimension: %w",
			err,
		)
	}

	if storedDimension != int64(ve.dimension) {
		return fmt.Errorf(
			"dimension mismatch: expected %d, got %d",
			ve.dimension,
			storedDimension,
		)
	}

	if count < 0 || count > maxVectors {
		return fmt.Errorf(
			"invalid vector count: %d",
			count,
		)
	}

	// Load IDs.
	for i := 0; i < int(count); i++ {
		var id int64

		if err := binary.Read(
			file,
			binary.LittleEndian,
			&id,
		); err != nil {
			return fmt.Errorf(
				"failed to read vector ID: %w",
				err,
			)
		}

		ve.ids[i] = int(id)
	}

	// Load vector data.
	vectorCount := int(count) * ve.dimension

	for i := 0; i < vectorCount; i++ {
		if err := binary.Read(
			file,
			binary.LittleEndian,
			&ve.vectorData[i],
		); err != nil {
			return fmt.Errorf(
				"failed to read vector data: %w",
				err,
			)
		}
	}

	// Only expose recovered data after the entire
	// snapshot has been successfully loaded.
	ve.count = int(count)

	return nil
}

// --------------------------------------------------
// WaitForPersistence
// --------------------------------------------------
//
// Useful for benchmarks/tests.
//
// It waits until all currently accepted persistence
// operations have completed.
//

func (ve *VectorEngine) WaitForPersistence() {
	ve.saveWG.Wait()
}
