package chat

type Message struct {
    ID        uint   `json:"id" gorm:"primaryKey"`
    Room      string `json:"room" gorm:"index"`
    User      string `json:"user"`
    Text      string `json:"text"`
    CreatedAt int64  `json:"created_at"`
}
