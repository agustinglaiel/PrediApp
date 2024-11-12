package repository

import (
	"context"
	"groupXUsers/internal/model"
	e "groupXUsers/pkg/utils"

	"gorm.io/gorm"
)

type groupXUserRepository struct {
	db *gorm.DB
}

type GroupXUserRepository interface {
	CreateGroupXUser(ctx context.Context, groupXUser *model.GroupXUser) e.ApiError
	GetGroupXUserByGroupID(ctx context.Context, groupId int) ([]*model.GroupXUser, e.ApiError)
	GetGroupXUserByGroupName(ctx context.Context, groupName string) ([]*model.GroupXUser, e.ApiError)
	GetGroupXUsers(ctx context.Context) ([]*model.GroupXUser, e.ApiError)
	RemoveUserFromGroup(ctx context.Context, userId int, groupId int) e.ApiError
	RemoveUserFromGroupByGroupName(ctx context.Context, userId int, groupName string) e.ApiError
}

func NewGroupXUserRepository(db *gorm.DB) GroupXUserRepository {
	return &groupXUserRepository{db: db}
}

// Para cuando un user se una a un group
func (r *groupXUserRepository) CreateGroupXUser(ctx context.Context, groupXUser *model.GroupXUser) e.ApiError {
	if err := r.db.WithContext(ctx).Create(&groupXUser).Error; err != nil {
		return e.NewInternalServerApiError("error creating groupXUser", err)
	}
	return nil
}

// Para obtener todos los usuarios que pertenecen a un GroupId
func (r *groupXUserRepository) GetGroupXUserByGroupID(ctx context.Context, groupId int) ([]*model.GroupXUser, e.ApiError) {
	var groupXUsers []*model.GroupXUser
	if err := r.db.WithContext(ctx).Where("group_id = ?", groupId).Find(&groupXUsers).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, e.NewNotFoundApiError("groupXUsers not found for group ID")
		}
		return nil, e.NewInternalServerApiError("error finding groupXUsers by group ID", err)
	}
	return groupXUsers, nil
}

// Para obtnener todos los usuarios que pertenecen a un GroupName
func (r *groupXUserRepository) GetGroupXUserByGroupName(ctx context.Context, groupName string) ([]*model.GroupXUser, e.ApiError) {
	var groupXUsers []*model.GroupXUser
	if err := r.db.WithContext(ctx).
		Joins("JOIN groups ON groups.id = group_x_users.group_id").
		Where("groups.group_name = ?", groupName).
		Find(&groupXUsers).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, e.NewNotFoundApiError("groupXUsers not found for group name")
		}
		return nil, e.NewInternalServerApiError("error finding groupXUsers by group name", err)
	}
	return groupXUsers, nil
}

// Obtiene todas las relaciones de la tabla
func (r *groupXUserRepository) GetGroupXUsers(ctx context.Context) ([]*model.GroupXUser, e.ApiError) {
	var groupXUsers []*model.GroupXUser
	if err := r.db.WithContext(ctx).Find(&groupXUsers).Error; err != nil {
		return nil, e.NewInternalServerApiError("error getting groupXUsers", err)
	}
	return groupXUsers, nil
}

func (r *groupXUserRepository) RemoveUserFromGroup(ctx context.Context, userId int, groupId int) e.ApiError {
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND group_id = ?", userId, groupId).
		Delete(&model.GroupXUser{}).Error; err != nil {
		return e.NewInternalServerApiError("error removing user from group", err)
	}
	return nil
}

func (r *groupXUserRepository) RemoveUserFromGroupByGroupName(ctx context.Context, userId int, groupName string) e.ApiError {
	if err := r.db.WithContext(ctx).
		Joins("JOIN groups ON groups.id = group_x_users.group_id").
		Where("group_x_users.user_id = ? AND groups.group_name = ?", userId, groupName).
		Delete(&model.GroupXUser{}).Error; err != nil {
		return e.NewInternalServerApiError("error removing user from group by name", err)
	}
	return nil
}
