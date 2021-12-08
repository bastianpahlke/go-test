package types

type HtmlData struct {
	UserId UserId `json:"userId"`
	Id     DataId `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type DataId int64
