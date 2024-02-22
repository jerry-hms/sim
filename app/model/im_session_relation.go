package model

import "strconv"

func CreateImSessionRelationFactory() *ImSessionRelation {
	return &ImSessionRelation{BaseModel: BaseModel{DB: ConnDb()}}
}

type ImSessionRelation struct {
	BaseModel
	UserId     uint64 `gorm:"column:user_id" json:"user_id"`
	RelationId uint64 `gorm:"column:relation_id" json:"relation_id"`
	Scene      string `gorm:"column:scene" json:"scene"`
	SepSvr     string `gorm:"column:sep_svr" json:"sep_svr"`
}

// GetRelationOrCreate 获取会话关系，如果不存在则创建
func (isr *ImSessionRelation) GetRelationOrCreate(user_id uint64, recv_id uint64, scene string) *ImSessionRelation {

	isr.UserId = user_id
	isr.RelationId = recv_id
	isr.Scene = scene
	isr.SepSvr = "0"

	isr.Where(
		"(user_id = ? and relation_id = ?) OR (relation_id = ? and user_id = ?)",
		user_id, recv_id, recv_id, user_id,
	).Where("scene = ?", scene).FirstOrCreate(&isr)

	return isr
}

func (isr *ImSessionRelation) IncSepSvr() {
	
	sepSvr, _ := strconv.Atoi(isr.SepSvr)
	isr.SepSvr = strconv.Itoa(sepSvr + 1)
	isr.Save(&isr)
}
