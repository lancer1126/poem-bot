package model

import (
	"go.uber.org/zap"
	"poem-bot/global"
	"strings"
)

const (
	BR = "\n"
)

//goland:noinspection SqlResolve
func GetRandomPoem() string {
	var poemMaps []map[string]interface{}

	sql := "select p.poem_name,e.era_name,a.author_name,p.anthology,p.chapter,p.section,p.poem_name,p.content from poem p " +
		" inner join era e on e.era_id = p.era_id " +
		" inner join author a on a.author_id = p.author_id " +
		" order by rand() limit 1"

	author := GetRandomAuthor()
	if author != "" {
		author = "'" + author + "'"
		sql = "select p.poem_name,e.era_name,a.author_name,p.anthology,p.chapter,p.section,p.poem_name,p.content from poem p " +
			" inner join era e on e.era_id = p.era_id " +
			" inner join author a on a.author_id = p.author_id " +
			" where a.author_name = " + author +
			" order by rand() limit 1"
	}

	err := global.DB.Raw(sql).Scan(&poemMaps).Error
	if err != nil {
		global.LOG.Error("failed to query database", zap.Error(err))
		return "从数据库获取错误"
	}

	if len(poemMaps) <= 0 {
		global.LOG.Error("content from database is empty")
		return "从数据库获取到空数据"
	}

	m := poemMaps[0]
	builder := strings.Builder{}
	builder.WriteString(toAppend(m, "poem_name"))
	builder.WriteString(toAppend(m, "chapter"))
	builder.WriteString(toAppend(m, "section"))
	builder.WriteString(toAppend(m, "author_name"))
	builder.WriteString(toAppend(m, "content"))
	return builder.String()
}

func GetRandomAuthor() string {
	var authorMaps []map[string]interface{}
	err := global.DB.Raw(`select author_name from match_author order by rand() limit 1`).Scan(&authorMaps).Error
	if err != nil {
		global.LOG.Error("failed to match author", zap.Error(err))
		return "从数据库获取诗人错误"
	}
	if len(authorMaps) <= 0 {
		return ""
	}
	return authorMaps[0]["author_name"].(string)
}

func toAppend(m map[string]interface{}, ele string) string {
	element := m[ele]
	if element == nil {
		return ""
	}
	if ele == "author_name" {
		return m["era_name"].(string) + "-" + m["author_name"].(string) + BR
	}
	return m[ele].(string) + BR
}
