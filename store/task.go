package store

import (
	"context"

	"github.com/yuyahy/go_todo_app/entity"
)

// taskの一覧を取得する
func (r *Repository) ListTasks(
	ctx context.Context, db Queryer,
) (entity.Tasks, error) {
	tasks := entity.Tasks{}
	sql := `SELECT id, title, status, created, modified FROM task;`

	if err := db.SelectContext(ctx, &tasks, sql); err != nil {
		return nil, err
	}
	return tasks, nil
}

// タスクを保存する
func (r *Repository) AddTask(
	ctx context.Context, db Execer, t *entity.Task,
) error {
	t.Created = r.Clocker.Now()
	t.Modified = r.Clocker.Now()
	sql := `INSERT INTO task
			(title, status, created, modified)
	VALUES (?, ?, ?, ?)`

	// sql := `INSERT INTO task (title, status, created, modified)
	// VALUES (?,?,?,?)`
	// クエリを実行
	result, err := db.ExecContext(
		ctx, sql, t.Title, t.Status, t.Created, t.Modified,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	// 引数のポインタにidをセット
	t.ID = entity.TaskID(id)
	return nil
}
