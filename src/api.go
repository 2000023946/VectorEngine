package src

type API struct {
	engine *VectorEngine
}

// NewAPI creates the VectorEngine and synchronously
// restores the persisted database from disk.
//
// The API is not returned until Reboot() completes.
// Therefore, callers cannot use the API until the
// previous persisted state has been loaded into RAM.
func NewAPI() (*API, error) {
	engine := NewVectorEngine()

	// Restore persisted state synchronously.
	//
	// If vectorengine.dat exists, Reboot() loads it
	// completely before NewAPI() returns.
	//
	// If it does not exist, the engine starts empty.
	if err := engine.Reboot(persistenceFile); err != nil {
		return nil, err
	}

	return &API{
		engine: engine,
	}, nil
}

// Insert adds a vector to the VectorEngine.
//
// VectorEngine handles:
//   - concurrent access
//   - in-memory storage
//   - persistence queue
//   - asynchronous disk persistence
//   - backpressure when the persistence queue is full
//
// Insert returns without waiting for the disk write
// itself to finish.
func (api *API) Insert(id int, values []float64) error {
	return api.engine.Insert(id, values)
}

// Search searches the VectorEngine for the K nearest vectors.
//
// Search operates entirely against the in-memory dataset.
func (api *API) Search(query []float64, k int) []Result {
	return api.engine.Search(query, k)
}

// Reset clears the database.
//
// Reset:
//  1. Prevents new inserts.
//  2. Clears the in-memory database.
//  3. Waits for all previously accepted persistence
//     operations to finish.
//  4. Deletes the persisted database from disk.
//  5. Allows inserts again.
//
// Searches during the reset simply see an empty
// in-memory database.
func (api *API) Reset() error {
	return api.engine.Reset()
}

// WaitForPersistence waits until all persistence
// operations currently accepted by the VectorEngine
// have completed.
//
// This is primarily useful for testing and benchmarking.
//
// Normal Insert() calls remain asynchronous and do not
// wait for persistence.
func (api *API) WaitForPersistence() {
	api.engine.WaitForPersistence()
}

// Count returns the number of vectors currently loaded in RAM.
func (api *API) Count() int {
	return api.engine.count
}
