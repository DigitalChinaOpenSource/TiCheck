package model

type ClusterChecklist struct {
	ID             uint       `gorm:"primarykey"`
	ClusterID      uint       `gorm:"not null"`
	ScriptID       uint       `gorm:"not null"`
	IsEnabled      int        `gorm:"not null"`
	Comparator     Comparator `gorm:"embedded"`
}

