package entity

import "time"

// 独自型を定義する事で誤代入を防ぐ
type TaskID int64
type TaskStatus string

const (
	TaskStatusTodo  TaskStatus = "todo"
	TaskStatusDoing TaskStatus = "doing"
	TaskStatusDone  TaskStatus = "done"
)

// tagを設定しておく事で構造体へマッピングするのが楽になる
// github.com/jmoiron/sqlx
type Task struct {
	ID       TaskID     `json:"id" db:"id"`
	UserID   UserID     `json:"user_id" db:"user_id"`
	Title    string     `json:"title" db:"title"`
	Status   TaskStatus `json:"status" db:"status"`
	Created  time.Time  `json:"created" db:"created"`
	Modified time.Time  `json:"modified" db:"modified"`
}

type Tasks []*Task
