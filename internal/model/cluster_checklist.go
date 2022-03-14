package model

type ClusterChecklist struct {
	ID         uint   `gorm:"primary" json:"id,omitempty"`
	ClusterID  uint   `gorm:"not null" json:"cluster_id,omitempty"`
	ProbeID    string `gorm:"not null" json:"probe_id,omitempty"`
	IsEnabled  int    `gorm:"not null" json:"is_enabled,omitempty"`
	Comparator `gorm:"embedded"`
}

type CheckListInfo struct {
	ID			uint   `gorm:"primary" json:"id"`
	ProbeID     string `gorm:"not null" json:"probe_id,omitempty"`
	ScriptName  string `gorm:"not null" json:"script_name,omitempty"`
	FileName    string `json:"file_name,omitempty"`
	Tag         string `json:"tag,omitempty"`
	Description string `json:"description,omitempty"`
	Operator    int    `json:"operator,omitempty"`
	Threshold   string `json:"threshold,omitempty"`
	IsEnabled   int    `json:"is_enabled,omitempty"`
}

func (cc *ClusterChecklist) TableName() string {
	return "cluster_checklist"
}

func (cc *ClusterChecklist) GetListInfoByClusterID(id int) ([]CheckListInfo, error) {
	var cl []CheckListInfo
	err := DbConn.Table("cluster_checklist as cc").Select("cc.id, cc.probe_id, p.script_name, p.file_name, " +
		"p.tag, p.description, cc.operator, cc.Threshold, cc.is_enabled").
		Joins("join probes as p on cc.probe_id = p.id").
		Where("cc.cluster_id = ?", id).Find(&cl).Error

	return cl, err
}

func (cc *ClusterChecklist) AddCheckProbe() error {
	err := DbConn.Create(cc).Error
	return err
}

func (cc *ClusterChecklist) DeleteCheckProbe() error {
	err := DbConn.Delete(cc).Error
	return err
}
