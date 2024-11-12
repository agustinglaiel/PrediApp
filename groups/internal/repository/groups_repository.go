package repository

import (
	"context"
	"groups/internal/model"
	e "groups/pkg/utils"

	"gorm.io/gorm"
)

type groupRespository struct {
	db *gorm.DB
}

type GroupRepository interface {
	CreateGroup(ctx context.Context, group *model.Group) e.ApiError
	GetGroupByID(ctx context.Context, id int) (*model.Group, e.ApiError)
	GetGroupByGroupName(ctx context.Context, groupName string) (*model.Group, e.ApiError)
	GetGroups(ctx context.Context) ([]*model.Group, e.ApiError)
	DeleteGroupByID(ctx context.Context, id int) e.ApiError
	DeleteGroupByGroupName(ctx context.Context, groupName string) e.ApiError
}

func NewGroupRepository(db *gorm.DB) GroupRepository {
	return &groupRespository{db: db}
}

func (r *groupRespository) CreateGroup(ctx context.Context, group *model.Group) e.ApiError {
	if err := r.db.WithContext(ctx).Create(&group).Error; err != nil {
		return e.NewInternalServerApiError("error creating group", err)
	}
	return nil
}

func (r *groupRespository) GetGroupByID(ctx context.Context, id int) (*model.Group, e.ApiError) {
	var group model.Group
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&group).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, e.NewNotFoundApiError("group not found")
		}
		return nil, e.NewInternalServerApiError("error finding group by id", err)
	}
	return &group, nil
}

func (r *groupRespository) GetGroupByGroupName(ctx context.Context, groupName string) (*model.Group, e.ApiError) {
	var group model.Group
	if err := r.db.WithContext(ctx).Where("group_name = ?", groupName).First(&group).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, e.NewNotFoundApiError("group not found")
		}
		return nil, e.NewInternalServerApiError("error finding group by group name", err)
	}
	return &group, nil
}

func (r *groupRespository) GetGroups(ctx context.Context) ([]*model.Group, e.ApiError) {
	var groups []*model.Group
	if err := r.db.WithContext(ctx).Find(&groups).Error; err != nil {
		return nil, e.NewInternalServerApiError("error getting groups", err)
	}
	return groups, nil
}

func (r *groupRespository) DeleteGroupByID(ctx context.Context, id int) e.ApiError {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Group{}).Error; err != nil {
		return e.NewInternalServerApiError("error deleting group", err)
	}
	return nil
}

func (r *groupRespository) DeleteGroupByGroupName(ctx context.Context, groupName string) e.ApiError {
	if err := r.db.WithContext(ctx).Where("group_name = ?", groupName).Delete(&model.Group{}).Error; err != nil {
		return e.NewInternalServerApiError("error deleting group", err)
	}
	return nil
}