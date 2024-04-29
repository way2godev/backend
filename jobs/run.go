package jobs

import "sync"

func RunJobs(started chan bool) {
	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)
	// go UpdateStopsCache()

	// Añadir aquí los jobs que se quieran ejecutar

	started <- true
}
