package model

import (
	"encoding/json"
)

const (
	TEXT MsgTypeType = "text"
)

type IDingMsg interface {
	Parse() []byte
}

type TextMsg struct {
	MsgType MsgTypeType `json:"msgtype,omitempty"`
	Text    textModel   `json:"text,omitempty"`
	At      atModel     `json:"at,omitempty"`
}

func (t TextMsg) Parse() []byte {
	b, _ := json.Marshal(t)
	return b
}

type textModel struct {
	Content string `json:"content,omitempty"`
}

type atModel struct {
	AtMobiles []string `json:"atMobiles,omitempty"`
	IsAtAll   bool     `json:"isAtAll,omitempty"`
}

type ResponseMsg struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type MsgTypeType string

func NewTextMsg(content string) *TextMsg {
	return &TextMsg{MsgType: TEXT, Text: textModel{Content: content}}
}
