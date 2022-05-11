package model

type ClusterChecklist struct {
	ID         uint   `gorm:"primary" json:"id"`
	ClusterID  uint   `gorm:"not null" json:"cluster_id"`
	ProbeID    string `gorm:"not null" json:"probe_id"`
	IsEnabled  int    `gorm:"not null" json:"is_enabled"`
	Comparator `gorm:"embedded"`
}

type CheckListInfo struct {
	ID          uint   `gorm:"primary" json:"id"`
	ProbeID     string `gorm:"not null" json:"probe_id"`
	ScriptName  string `gorm:"not null" json:"script_name"`
	FileName    string `json:"file_name"`
	Tag         string `json:"tag"`
	Source      string `json:"source"` // source is one of ["local","remote","custom"]
	Description string `json:"description"`
	Operator    int    `json:"operator"`
	Threshold   string `json:"threshold"`
	IsEnabled   int    `json:"is_enabled"`
	Arg         string `json:"arg,omitempty"`
}

func (cc *ClusterChecklist) TableName() string {
	return "cluster_checklist"
}

func (cc *ClusterChecklist) GetListInfoByClusterID(id int) ([]CheckListInfo, error) {
	var cl []CheckListInfo
	var probe Probe
	err := DbConn.Table(cc.TableName()+" as cc").Select("cc.id, cc.probe_id, p.script_name, p.file_name, "+
		"p.tag, case p.is_system when 1 then 'local' when 0 then 'custom' else 'remote' end as source, p.description, cc.operator, cc.threshold, cc.is_enabled, cc.arg").
		Joins("join "+probe.TableName()+" as p on cc.probe_id = p.id").
		Where("cc.cluster_id = ?", id).Find(&cl).Error

	return cl, err
}

func (cc *ClusterChecklist) GetEnabledCheckListByClusterID(id int) ([]CheckListInfo, error) {
	var cl []CheckListInfo
	var probe Probe
	err := DbConn.Table(cc.TableName()+" as cc").Select("cc.id, cc.probe_id, p.script_name, p.file_name, "+
		"p.tag, case p.is_system when 1 then 'local' when 0 then 'custom' else 'remote' end as source, p.description, cc.operator, cc.threshold, cc.is_enabled, cc.arg").
		Joins("join "+probe.TableName()+" as p on cc.probe_id = p.id").
		Where("cc.is_enabled = 1 and cc.cluster_id = ?", id).Find(&cl).Error

	return cl, err
}

func (cc *ClusterChecklist) GetEnabledCheckListTagGroup(id int) map[string]interface{} {
	result := make(map[string]interface{})
	var total int

	groups := make(map[string]interface{})
	var probe Probe
	rows, _ := DbConn.Table(cc.TableName()+" as cc").Select("p.tag, count(cc.probe_id) as cnt").
		Joins("join "+probe.TableName()+" as p on cc.probe_id = p.id").
		Where("cc.is_enabled = 1 and cc.cluster_id = ?", id).Group("p.tag").Rows()
	defer rows.Close()

	for rows.Next() {
		g := struct {
			Tag string
			Cnt int
		}{}
		DbConn.ScanRows(rows, &g)
		total += g.Cnt
		groups[g.Tag] = g.Cnt
	}

	result["tags"] = groups
	result["total"] = total
	return result
}

func (cc *ClusterChecklist) AddCheckProbe() error {
	err := DbConn.Create(cc).Error
	return err
}

func (cc *ClusterChecklist) AddCheckList(checkList []ClusterChecklist) error {
	err := DbConn.Create(checkList).Error
	return err
}

func (cc *ClusterChecklist) DeleteCheckProbe() error {
	err := DbConn.Delete(cc).Error
	return err
}

// ChangeProbeStatus Only update status for is_enabled by id
func (cc *ClusterChecklist) ChangeProbeStatus() error {
	err := DbConn.Model(cc).Update("is_enabled", cc.IsEnabled).Error
	return err
}

// UpdateProbeConfig only update operator and threshold by id
func (cc *ClusterChecklist) UpdateProbeConfig() error {
	err := DbConn.Model(cc).Updates(map[string]interface{}{
		"operator":  cc.Operator,
		"threshold": cc.Threshold,
	}).Error
	return err
}
