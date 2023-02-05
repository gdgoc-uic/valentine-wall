package main

type UserConnection struct {
	UserID      string `db:"user_id" json:"-"`
	Provider    string `db:"provider" json:"provider"`
	Token       string `db:"token" json:"-"`
	TokenSecret string `db:"token_secret" json:"-"`
}

type RecipientStats struct {
	RecipientID       string `db:"recipient_id" json:"recipient_id"`
	MessagesCount     int    `db:"messages_count" json:"messages_count"`
	GiftMessagesCount int    `db:"gift_messages_count" json:"gift_messages_count"`
}

type RecipientStats2 struct {
	RecipientID string  `db:"recipient_id" json:"recipient_id"`
	Department  string  `db:"department" json:"department"`
	Sex         string  `db:"sex" json:"sex"`
	TotalCoins  float32 `db:"-" json:"total_coins"`
}

type Recipients []*RecipientStats2

func (a Recipients) Len() int           { return len(a) }
func (a Recipients) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Recipients) Less(i, j int) bool { return a[i].TotalCoins > a[j].TotalCoins }

func (a Recipients) BySex(sex string) Recipients {
	if len(sex) == 0 || sex == "all" {
		return a
	}
	res := make(Recipients, 0, len(a))
	for _, v := range a {
		if v.Sex == sex {
			res = append(res, v)
		}
	}
	return res
}
