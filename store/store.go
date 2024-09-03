package store

import (
	"errors"

	"github.com/yuyahy/go_todo_app/entity"
)

var (
	Tasks = &TaskStore{Tasks: map[entity.TaskID]*entity.Task{}}

	ErrNotFound = errors.New("Not found")
)

type TaskStore struct {
	// 動作確認の仮実装なのであえてexportしている
	LastID entity.TaskID
	Tasks  map[entity.TaskID]*entity.Task
}

func (ts *TaskStore) Add(t *entity.Task) (entity.TaskID, error) {
	ts.LastID++
	t.ID = ts.LastID
	ts.Tasks[t.ID] = t
	return t.ID, nil
}

// ソート済みのタスクを返す
func (ts *TaskStore) All() entity.Tasks {
	tasks := make([]*entity.Task, len(ts.Tasks))
	// TaskStore.Tasksのmapのkey,valueを取り出してループ
	for i, t := range ts.Tasks {
		tasks[i-1] = t
	}
	return tasks
}
