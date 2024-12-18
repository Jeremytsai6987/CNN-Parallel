package scheduler

import (
	"fmt"
)


func RunParallelSlices(config Config) {
	effects, _ := ReadEffectsFile()
	
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
	for i := 0; i < numTasks; i++ {
		task := queue.Dequeue()
		ApplyEffectsToSlice(task.InPath, task.OutPath, task.Effects, numThreads)
	}

}
