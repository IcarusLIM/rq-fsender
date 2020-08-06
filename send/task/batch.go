package task

import (
	"github.com/Ghamster0/os-rq-fsender/pkg/dto"
	"github.com/Ghamster0/os-rq-fsender/pkg/sth"
	"github.com/Ghamster0/os-rq-fsender/send/entity"
	"github.com/jinzhu/gorm"
)

type BatchService struct {
	db  *gorm.DB
	box *TaskBox
}

func NewBatchService(db *gorm.DB, box *TaskBox) *BatchService {
	return &BatchService{
		db:  db,
		box: box,
	}
}

func (bs *BatchService) CreateBatch(batchReq *dto.BatchReq) (sth.Result, error) {
	var add, ignore []string
	var totalSize int64
	batchModel := &entity.BatchModel{
		Id:  batchReq.Id,
		Api: batchReq.Api,
	}
	err := bs.db.Transaction(
		func(tx *gorm.DB) error {
			if err := tx.Create(batchModel).Error; err != nil {
				return err
			}
			for i := range batchReq.FileIds {
				fid := batchReq.FileIds[i]
				var fm entity.FileModel
				if err := tx.Where("id = ? AND status = ?", fid, entity.Idel).First(&fm).Error; err != nil {
					ignore = append(ignore, fid)
					continue
				}
				if err := tx.Model(&fm).Where("id = ?", fid).Update(sth.Map{"status": entity.Waitting, "batch_id": batchReq.Id}).Error; err != nil {
					ignore = append(ignore, fid)
					continue
				}
				add = append(add, fid)
				totalSize += fm.Size
			}
			tx.Model(&batchModel).Where("id = ?", batchModel.Id).Update(sth.Map{"total_size": totalSize})
			return nil
		},
	)
	res := sth.Result{
		"batch":  &batchModel,
		"add":    add,
		"ignore": ignore,
	}
	return res, err
}

func (bs *BatchService) ListBatch(start int, limit int) (batches []entity.BatchModel, err error) {
	err = bs.db.Offset(start).Limit(limit).Find(&batches).Error
	return
}

func (bs *BatchService) GetBatch(id string) (sth.Result, error) {
	var model entity.BatchModel
	if err := bs.db.Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	var fms []entity.FileModel
	if err := bs.db.Where("batch_id = ?", id).Find(&fms).Error; err != nil {
		return nil, err
	}
	var taskInfos []sth.Result
	for i := range fms {
		fm := fms[i]
		info, err := bs.box.InfoTask(fm.Id)
		if err != nil {
			info, err = fm.Info()
		}
		if err == nil {
			taskInfos = append(taskInfos, info)
		}
	}
	return sth.Result{
		"id":         model.Id,
		"api":        model.Api,
		"total_size": model.TotalSize,
		"infos":      taskInfos,
	}, nil
}

func (bs *BatchService) DelBatch(id string) (sth.Result, error) {
	_, fms, err := bs.GetBatchFiles(id)
	if err != nil {
		return nil, err
	}
	var taskInfos []sth.Result
	for _, fm := range *fms {
		info, err := bs.box.CancelTask(fm.Id)
		if err != nil {
			info, err = fm.Info()
		}
		if err == nil {
			taskInfos = append(taskInfos, info)
		}
	}
	bs.DelBatchFiles(id)
	return sth.Result{
		"infos": taskInfos,
	}, nil
}

func (bs *BatchService) GetBatchFiles(id string) (*entity.BatchModel, *[]entity.FileModel, error) {
	var model entity.BatchModel
	if err := bs.db.Where("id = ?", id).First(&model).Error; err != nil {
		return nil, nil, err
	}
	var fms []entity.FileModel
	if err := bs.db.Where("batch_id = ?", id).Find(&fms).Error; err != nil {
		return &model, nil, err
	}
	return &model, &fms, nil
}

func (bs *BatchService) DelBatchFiles(id string) error {
	err := bs.db.Transaction(
		func(tx *gorm.DB) error {
			if err := tx.Where("id = ?", id).Delete(entity.BatchModel{}).Error; err != nil {
				return err
			}
			if err := tx.Where("batch_id = ?", id).Delete(entity.FileModel{}).Error; err != nil {
				return err
			}
			return nil
		})
	return err
}
