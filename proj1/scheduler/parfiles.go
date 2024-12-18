package scheduler

import (
	"fmt"
	"sync"
)
func worker(id int, queue *TaskQueue, lock *TASLock, wg *sync.WaitGroup){
	defer wg.Done()
	for{
		lock.Lock()
		task := queue.Dequeue()
		lock.Unlock()
		if task == nil{
			return
		}
		err := ApplyEffects(task.InPath, task.OutPath, task.Effects)
	}
}

func RunParallelFiles(config Config) {
	effects, err := ReadEffectsFile()
	if err != nil{
		fmt.Println("Error reading effects file")
		return
	}
	queue := TaskQueue{}
	for _, effect := range effects{
        task := &ImageTask{
            InPath:  fmt.Sprintf("../data/in/%s/%s", config.DataDirs, effect.InPath),
            OutPath: fmt.Sprintf("../data/out/%s_%s", config.DataDirs, effect.OutPath),
            Effects: effect.Effects,
        }
        queue.Enqueue(task)	
	}
	numThreads := config.ThreadCount
	numTasks := len(effects)
	numWorkers := numThreads 
    if numWorkers > numTasks {
        numWorkers = numTasks
    }	
	lock := &TASLock{}
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++{
		go worker(i, &queue, lock, &wg)
	}
	wg.Wait()
}
