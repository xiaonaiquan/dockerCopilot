package utiles

import (
	"sync"
	"time"
)

type TaskStatus string

const (
	TaskPending TaskStatus = "pending"
	TaskRunning TaskStatus = "running"
	TaskSuccess TaskStatus = "success"
	TaskFailed  TaskStatus = "failed"
)

type Task struct {
	ID      string        // 任务ID
	Status  TaskStatus    // 状态
	Log     []string      // 日志/进度
	Created time.Time     // 创建时间
	Done    chan struct{} // 用于通知任务完成
}

type TaskManager struct {
	tasks map[string]*Task
	mu    sync.RWMutex
}

var GlobalTaskManager = NewTaskManager()

func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks: make(map[string]*Task),
	}
}

func (tm *TaskManager) AddTask(task *Task) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.tasks[task.ID] = task
}

func (tm *TaskManager) GetTask(id string) (*Task, bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	t, ok := tm.tasks[id]
	return t, ok
}

func (tm *TaskManager) UpdateTaskLog(id string, log string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	if t, ok := tm.tasks[id]; ok {
		t.Log = append(t.Log, log)
	}
}

func (tm *TaskManager) UpdateTaskStatus(id string, status TaskStatus) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	if t, ok := tm.tasks[id]; ok {
		t.Status = status
		if status == TaskSuccess || status == TaskFailed {
			close(t.Done)
		}
	}
}
