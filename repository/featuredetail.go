package repository

import (
	"documentation/entity"
	"documentation/input"
	"fmt"

	"gorm.io/gorm"
)

type FeatureDetailRepository interface {
	SaveFeatureDetail(featuredetail entity.FeatureDetail) (entity.FeatureDetail, error)
	UpdateFeatureDetail(featuredetail entity.FeatureDetail) (entity.FeatureDetail, error)
	FindByIDFeatureDetail(ID int) (entity.FeatureDetail, error)
	FindAllFeatureDetail(inputID input.InputIDFeatureDetail) ([]entity.FeatureDetail, error)
	DeleteByIDFeatureDetail(ID int) (entity.FeatureDetail, error)
	FindAll(search string, page int, size int) ([]entity.FeatureDetail, error)
	TotalFetchData(search string, page int, size int) (int, error)
}

type featuredetailRepository struct {
	db *gorm.DB
}

func NewFeatureDetailRepository(db *gorm.DB) *featuredetailRepository {
	return &featuredetailRepository{db}
}

func (r *featuredetailRepository) SaveFeatureDetail(featuredetail entity.FeatureDetail) (entity.FeatureDetail, error) {
	err := r.db.Create(&featuredetail).Error
	if err != nil {
		return featuredetail, err
	}
	return featuredetail, nil

}
func (r *featuredetailRepository) FindByIDFeatureDetail(ID int) (entity.FeatureDetail, error) {
	var featuredetail entity.FeatureDetail
	err := r.db.Where("id = ? ", ID).Find(&featuredetail).Error
	if err != nil {
		return featuredetail, err
	}
	return featuredetail, nil

}
func (r *featuredetailRepository) UpdateFeatureDetail(featuredetail entity.FeatureDetail) (entity.FeatureDetail, error) {
	err := r.db.Save(&featuredetail).Error
	if err != nil {
		return featuredetail, err
	}
	return featuredetail, nil

}
func (r *featuredetailRepository) FindAllFeatureDetail(inputID input.InputIDFeatureDetail) ([]entity.FeatureDetail, error) {
	var featuredetails []entity.FeatureDetail
	err := r.db.Where("feature_id = ? ", inputID.ID).Find(&featuredetails).Error
	if err != nil {
		return featuredetails, err
	}
	return featuredetails, nil

}
func (r *featuredetailRepository) DeleteByIDFeatureDetail(ID int) (entity.FeatureDetail, error) {
	var featuredetail entity.FeatureDetail
	err := r.db.Where("id = ? ", ID).Delete(&featuredetail).Error
	if err != nil {
		return featuredetail, err
	}
	return featuredetail, nil

}

func (r *featuredetailRepository) FindAll(search string, page int, size int) ([]entity.FeatureDetail, error) {
	var feature []entity.FeatureDetail
	searchQuery := search
	sql := "SELECT a.id,a.feature_id,b.feature_name,path_of_file,a.title,a.code,a.description ,a.author_name ,a.created_at ,a.updated_at  FROM feature_details a join features b on a.feature_id=b.id WHERE 1=1  "
	if searchQuery != "" {
		sql = fmt.Sprintf("%s AND (a.feature_name LIKE '%%%s%%' OR a.title LIKE '%%%s%%') ", sql, searchQuery, searchQuery)
	}
	start := (page * size)
	limit := 0
	if start < 0 {
		limit = 0
	} else {
		limit = start
	}
	length := size

	sql = fmt.Sprintf("%s ORDER BY a.id DESC LIMIT %d OFFSET %d", sql, length, limit)

	err := r.db.Raw(sql).Scan(&feature).Error
	if err != nil {
		return feature, err
	}

	return feature, nil

}

func (r *featuredetailRepository) TotalFetchData(search string, page int, size int) (int, error) {
	var feature []entity.FeatureDetail
	searchQuery := search
	sql := "SELECT a.id,a.feature_id,b.feature_name,path_of_file,a.title,a.code,a.description ,a.author_name ,a.created_at ,a.updated_at  FROM feature_details a join features b on a.feature_id=b.id WHERE 1=1 "
	if searchQuery != "" {
		sql = fmt.Sprintf("%s AND (a.feature_name LIKE '%%%s%%' OR a.title LIKE '%%%s%%') ", sql, searchQuery, searchQuery)
	}

	err := r.db.Raw(sql).Scan(&feature).Error
	if err != nil {
		return len(feature), err
	}
	return len(feature), nil
}

//Generated by Micagen at 25 Mei 2022