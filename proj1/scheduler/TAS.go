package scheduler

import (
	"runtime"
	"sync/atomic"
)

type ImageTask struct {
	InPath string
	OutPath string
	Effects []string
	next *ImageTask
}

type TaskQueue struct {
	head *ImageTask
	tail *ImageTask
}

func (q *TaskQueue) Enqueue(task *ImageTask) {
	if q.head == nil {
		q.head = task
		q.tail = task
	} else {
		q.tail.next = task
		q.tail = task
	}
}
// FIFO queue
func (q *TaskQueue) Dequeue() (task *ImageTask) {
	if q.head == nil {
		return nil
	}
	task = q.head
	q.head = q.head.next
	return task
}

type TASLock struct {
	locked int32
}

func (lock *TASLock) Lock(){
	for !atomic.CompareAndSwapInt32(&lock.locked, 0, 1){
		runtime.Gosched()
	}
}

func (lock *TASLock) Unlock(){
	atomic.StoreInt32(&lock.locked, 0)
}
