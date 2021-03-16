package main

type updatesModel struct {
	Ok     bool          `json:"ok"`
	Result []resultModel `json:"result"`
}

type resultModel struct {
	UpdateId int64        `json:"update_id"`
	Message  messageModel `json:"message"`
}

type messageModel struct {
	MessageId int64           `json:"message_id"`
	Date      int64           `json:"date"`
	Text      string          `json:"text"`
	From      fromModel       `json:"from"`
	Chat      chatModel       `json:"chat"`
	Entities  []entitiesModel `json:"entities"`
}

type fromModel struct {
	Id        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
}

type chatModel struct {
	Id   int64  `json:"id"`
	Type string `json:"Type"`
}

type entitiesModel struct {
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	Type   string `json:"type"`
}
