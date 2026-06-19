//=============================================================================
//===
//=== Copyright (C) 2024-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package core

import (
	"log/slog"
	"time"
)

//=============================================================================

type WorkerPool struct {
	numWorkers int
	queueSize  int
	taskQueue  chan func()
}

//=============================================================================
//===
//=== Methods
//===
//=============================================================================

func (p *WorkerPool) Init(numWorkers, queueSize int) {
	p.numWorkers = numWorkers
	p.queueSize  = queueSize
	p.taskQueue  = make(chan func(), queueSize)

	for i := 0; i<numWorkers; i++ {
		go p.worker()
	}

	slog.Info("Init: Worker pool created", "workers", numWorkers)
}

//=============================================================================

func (p *WorkerPool) ShutDown() bool {

	if p.taskQueue != nil {
		close(p.taskQueue)

		//--- Wait for task completion
		for ; len(p.taskQueue) > 0 ; {
			time.Sleep(time.Millisecond * 500)
		}

		p.taskQueue = nil
	}

	slog.Info("ShutDown: Worker pool stopped")

	return true
}

//=============================================================================

func (p *WorkerPool) Submit(task func()) {
	p.taskQueue <- task
}

//=============================================================================
//===
//=== Worker
//===
//=============================================================================

func (p *WorkerPool) worker() {
	for {
		select {
			case task, ok := <- p.taskQueue:
				if ok {
					//--- Run task
					task()
				} else {
					//--- Channel closed. Exit from goroutine
					return
				}

			default:
				time.Sleep(time.Millisecond * 50)
		}
	}
}

//=============================================================================
