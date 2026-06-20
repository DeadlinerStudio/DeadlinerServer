package document

type SubTask struct {
	ID          string `json:"id"`
	Content     string `json:"content"`
	IsCompleted bool   `json:"is_completed"`
	SortOrder   int32  `json:"sort_order"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
