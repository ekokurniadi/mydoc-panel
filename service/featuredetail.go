package service

import (
	"documentation/entity"
	"documentation/input"
	"documentation/repository"
)

type FeatureDetailService interface {
	FeatureDetailServiceGetAll(inputID input.InputIDFeatureDetail) ([]entity.FeatureDetail, error)
	FeatureDetailServiceGetByID(inputID input.InputIDFeatureDetail) (entity.FeatureDetail, error)
	FeatureDetailServiceCreate(input input.FeatureDetailInput) (entity.FeatureDetail, error)
	FeatureDetailServiceUpdate(inputID input.InputIDFeatureDetail, inputData input.FeatureDetailInput) (entity.FeatureDetail, error)
	FeatureDetailServiceDeleteByID(inputID input.InputIDFeatureDetail) (bool, error)
	FindAll(search string, page int, size int) ([]entity.FeatureDetail, error)
	TotalFetchData(search string, page int, size int) (int, error)
}
type featuredetailService struct {
	repository repository.FeatureDetailRepository
}

func NewFeatureDetailService(repository repository.FeatureDetailRepository) *featuredetailService {
	return &featuredetailService{repository}
}
func (s *featuredetailService) FeatureDetailServiceCreate(input input.FeatureDetailInput) (entity.FeatureDetail, error) {
	featuredetail := entity.FeatureDetail{}
	featuredetail.FeatureID = input.FeatureID
	featuredetail.PathOfFile = input.PathOfFile
	featuredetail.Title = input.Title
	featuredetail.Code = input.Code
	featuredetail.Description = input.Description
	featuredetail.AuthorName = input.AuthorName
	newFeatureDetail, err := s.repository.SaveFeatureDetail(featuredetail)
	if err != nil {
		return newFeatureDetail, err
	}
	return newFeatureDetail, nil
}
func (s *featuredetailService) FeatureDetailServiceUpdate(inputID input.InputIDFeatureDetail, inputData input.FeatureDetailInput) (entity.FeatureDetail, error) {
	featuredetail, err := s.repository.FindByIDFeatureDetail(inputID.ID)
	if err != nil {
		return featuredetail, err
	}
	featuredetail.FeatureID = inputData.FeatureID
	featuredetail.PathOfFile = inputData.PathOfFile
	featuredetail.Title = inputData.Title
	featuredetail.Code = inputData.Code
	featuredetail.Description = inputData.Description
	featuredetail.AuthorName = inputData.AuthorName

	updatedFeatureDetail, err := s.repository.UpdateFeatureDetail(featuredetail)

	if err != nil {
		return updatedFeatureDetail, err
	}
	return updatedFeatureDetail, nil
}
func (s *featuredetailService) FeatureDetailServiceGetByID(inputID input.InputIDFeatureDetail) (entity.FeatureDetail, error) {
	featuredetail, err := s.repository.FindByIDFeatureDetail(inputID.ID)
	if err != nil {
		return featuredetail, err
	}
	return featuredetail, nil
}
func (s *featuredetailService) FeatureDetailServiceGetAll(inputID input.InputIDFeatureDetail) ([]entity.FeatureDetail, error) {
	featuredetails, err := s.repository.FindAllFeatureDetail(inputID)
	if err != nil {
		return featuredetails, err
	}
	return featuredetails, nil
}
func (s *featuredetailService) FeatureDetailServiceDeleteByID(inputID input.InputIDFeatureDetail) (bool, error) {
	_, err := s.repository.FindByIDFeatureDetail(inputID.ID)
	if err != nil {
		return false, err
	}
	_, err = s.repository.DeleteByIDFeatureDetail(inputID.ID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *featuredetailService) FindAll(search string, page int, size int) ([]entity.FeatureDetail, error) {
	features, err := s.repository.FindAll(search, page, size)
	if err != nil {
		return features, err
	}
	return features, nil
}

func (s *featuredetailService) TotalFetchData(search string, page int, size int) (int, error) {
	totalData, err := s.repository.TotalFetchData(search, page, size)
	if err != nil {
		return totalData, err
	}
	return totalData, nil
}

//Generated by Micagen at 25 Mei 2022