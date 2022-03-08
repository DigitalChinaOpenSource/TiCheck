package model

type Script struct {
	ID                uint   `gorm:"primarykey"`
	ScriptName        string `gorm:"not null;unique"`
	FileName          string `gorm:"not null;unique"`
	Tag               string `gorm:"not null"`
	Description       string
	DefaultComparator Comparator `gorm:"embedded"`
	IsSystem          int        `gorm:"not null"`
	Creator           string     `gorm:"not null"`
	CreateTime        JsonTime   `gorm:"not null"`
	UpdateTime        JsonTime   `gorm:"not null"`
}

type ScriptMeta struct {
	ID          string        `json:"_id"`
	Name        string        `json:"name"`
	Author      ScriptAuthor  `json:"author"`
	Description string        `json:"description"`
	Tags        []string      `json:"tags"`
	Comparators []*Comparator `json:"rules"`
	Files       []string      `json:"files"`
	Main        string        `json:"main"`
	HomePage    string        `json:"homepage"`
	Version     string        `json:"version"`
	CreateTime  JsonTime      `json:"createTime"`
	UpdateTime  JsonTime      `json:"updateTime"`
}

type ScriptAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
