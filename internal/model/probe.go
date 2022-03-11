package model

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Probe struct {
	ID          string `gorm:"primarykey" json:"id"`
	ScriptName  string `gorm:"not null;unique" json:"script_name"`
	FileName    string `gorm:"not null;unique" json:"file_name"`
	Tag         string `gorm:"not null" json:"tag"`
	Description string `json:"description"`
	Comparator  `gorm:"embedded"`
	IsSystem    int       `gorm:"not null" json:"-"`
	Creator     string    `gorm:"not null" json:"creator"`
	CreateTime  time.Time `gorm:"not null" json:"create_time"`
	UpdateTime  time.Time `gorm:"not null" json:"update_time"`
}

func (Probe) TableName() string {
	return "probes"
}

type ProbeMeta struct {
	ID          string        `json:"_id"`
	Name        string        `json:"name"`
	Author      ProbeAuthor   `json:"author"`
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

type ProbeAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (p *Probe) GetPager(c *gin.Context, pg *Paginator) *Paginator {
	pg.SetPager(c.Request, "update_time desc")

	var rows []Probe
	return pg.ApplyQuery(DbConn.Model(&Probe{}), &rows)
}

func (p Probe) GetNotAddedProveListByClusterID(id int) ([]Probe, error) {
	var probeId []string
	var probeList []Probe

	// First, finding all installed probes in cluster
	err := DbConn.Table("cluster_checklist").Select("cluster_checklist.probe_id").Where("cluster_id = ?", id).Find(&probeId).Error
	if err !=  nil {
		return probeList, err
	}

	if len(probeId) > 0 {
		err = DbConn.Where("id not in ?", probeId).Find(&probeList).Error
	} else {
		err = DbConn.Find(&probeList).Error
	}

	return probeList, err
}