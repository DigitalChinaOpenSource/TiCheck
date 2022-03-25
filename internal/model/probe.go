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
	IsSystem    int       `gorm:"not null" json:"is_system"`
	Creator     string    `gorm:"not null" json:"creator"`
	CreateTime  time.Time `gorm:"not null" json:"create_time"`
	UpdateTime  time.Time `gorm:"not null" json:"update_time"`
}

func (Probe) TableName() string {
	return "probes"
}

type ProbeMeta struct {
	ID          string          `json:"_id"`
	Name        string          `json:"name"`
	Author      ProbeAuthor     `json:"author"`
	Description string          `json:"description"`
	Tags        []string        `json:"tags"`
	Rules       []ProbeMetaRule `json:"rules"`
	Files       []string        `json:"files"`
	Main        string          `json:"main"`
	HomePage    string          `json:"homepage"`
	Version     string          `json:"version"`
	CreateTime  JsonTime        `json:"createTime"`
	UpdateTime  JsonTime        `json:"updateTime"`
}

type ProbeMetaRule struct {
	// represents way to compare threshold and result
	// 0: NA, 1: eq. 2: g, 3: ge, 4: l, 5: le
	Operator  int      `json:"operator"`
	Threshold string   `json:"threshold"`
	Args      []string `json:"args,omitempty"`
}

type ProbeAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (p *Probe) Create() error {
	return DbConn.Create(p).Error
}

func (p *Probe) GetPager(c *gin.Context, pg *Paginator) *Paginator {
	pg.SetPager(c.Request, "update_time desc")

	var rows []Probe
	return pg.ApplyQuery(DbConn.Model(&Probe{}), &rows)
}

func (p *Probe) GetByID() error {
	return DbConn.Where("id = ?", p.ID).First(p).Error
}

func (p *Probe) Delete() error {
	return DbConn.Delete(p).Error
}

func (p *Probe) IsNotExist() bool {
	var cnt int64
	DbConn.Table(p.TableName()).Where("id = ?", p.ID).Count(&cnt)
	return cnt == 0
}

func (p Probe) GetNotAddedProveListByClusterID(id int) ([]Probe, error) {
	var probeId []string
	var probeList []Probe
	var cc ClusterChecklist

	// First, finding all installed probes in cluster
	err := DbConn.Table(cc.TableName()).Select("probe_id").Where("cluster_id = ?", id).Find(&probeId).Error
	if err != nil {
		return probeList, err
	}

	if len(probeId) > 0 {
		err = DbConn.Where("id not in ?", probeId).Find(&probeList).Error
	} else {
		err = DbConn.Find(&probeList).Error
	}

	return probeList, err
}
