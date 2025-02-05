package repository

import (
	"context"
	"groups/internal/model"
	e "groups/pkg/utils"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

type groupRepository struct {
	db *gorm.DB
}

type GroupRepository interface {
	CreateGroup(ctx context.Context, group *model.Group) e.ApiError
	GetGroupByID(ctx context.Context, id int) (*model.Group, e.ApiError)
	GetGroupByGroupName(ctx context.Context, groupName string) (*model.Group, e.ApiError)
	GetGroupByCode(ctx context.Context, code string) (*model.Group, e.ApiError)  
	GetGroups(ctx context.Context) ([]*model.Group, e.ApiError)
	DeleteGroupByID(ctx context.Context, id int) e.ApiError
	DeleteGroupByGroupName(ctx context.Context, groupName string) e.ApiError

	// Manejo de `groups_users`
	AddUserToGroup(ctx context.Context, groupID int, userID int, role string) e.ApiError 
	RemoveUserFromGroup(ctx context.Context, groupID int, userID int) e.ApiError         
	UserExistsInGroup(ctx context.Context, groupID int, userID int) (bool, e.ApiError)   
}

func NewGroupRepository(db *gorm.DB) GroupRepository {
	return &groupRepository{db: db}
}

func (r *groupRepository) CreateGroup(ctx context.Context, group *model.Group) e.ApiError {
	for {
		group.GroupCode = generateGroupCode() // Generar código aleatorio

		// Verificar si el código ya existe
		existingGroup, _ := r.GetGroupByCode(ctx, group.GroupCode)
		if existingGroup == nil {
			break // Si no existe, usamos este código
		}
	}

	if err := r.db.WithContext(ctx).Create(&group).Error; err != nil {
		return e.NewInternalServerApiError("Error creating group", err)
	}
	return nil
}

func (r *groupRepository) GetGroupByID(ctx context.Context, id int) (*model.Group, e.ApiError) {
	var group model.Group
	if err := r.db.WithContext(ctx).
		Preload("GroupUsers"). // Cargamos correctamente los usuarios del grupo
		Where("id = ?", id).
		First(&group).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, e.NewNotFoundApiError("Group not found")
		}
		return nil, e.NewInternalServerApiError("Error finding group by ID", err)
	}
	return &group, nil
}

func (r *groupRepository) GetGroupByGroupName(ctx context.Context, groupName string) (*model.Group, e.ApiError) {
	var group model.Group
	if err := r.db.WithContext(ctx).Where("group_name = ?", groupName).First(&group).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, e.NewNotFoundApiError("group not found")
		}
		return nil, e.NewInternalServerApiError("error finding group by group name", err)
	}
	return &group, nil
}

func (r *groupRepository) GetGroupByCode(ctx context.Context, code string) (*model.Group, e.ApiError) {
	var group model.Group
	if err := r.db.WithContext(ctx).
		Preload("Users").
		Where("group_code = ?", code).
		First(&group).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, e.NewNotFoundApiError("Group not found")
		}
		return nil, e.NewInternalServerApiError("Error finding group by code", err)
	}
	return &group, nil
}

func (r *groupRepository) GetGroups(ctx context.Context) ([]*model.Group, e.ApiError) {
	var groups []*model.Group
	if err := r.db.WithContext(ctx).
		Preload("Users").
		Find(&groups).Error; err != nil {
		return nil, e.NewInternalServerApiError("Error getting groups", err)
	}
	return groups, nil
}

func (r *groupRepository) DeleteGroupByID(ctx context.Context, id int) e.ApiError {
	// Eliminar usuarios del grupo primero
	if err := r.db.WithContext(ctx).
		Where("group_id = ?", id).
		Delete(&model.GroupXUsers{}).Error; err != nil {
		return e.NewInternalServerApiError("Error deleting users from group", err)
	}

	// Luego eliminar el grupo
	if err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&model.Group{}).Error; err != nil {
		return e.NewInternalServerApiError("Error deleting group", err)
	}
	return nil
}

func (r *groupRepository) DeleteGroupByGroupName(ctx context.Context, groupName string) e.ApiError {
	if err := r.db.WithContext(ctx).Where("group_name = ?", groupName).Delete(&model.Group{}).Error; err != nil {
		return e.NewInternalServerApiError("Error deleting group", err)
	}
	return nil
}

func (r *groupRepository) AddUserToGroup(ctx context.Context, groupID int, userID int, role string) e.ApiError {
	exists, err := r.UserExistsInGroup(ctx, groupID, userID)
	if err != nil {
		return err
	}
	if exists {
		return e.NewBadRequestApiError("User is already in the group")
	}

	groupUser := model.GroupXUsers{
		GroupID:   groupID,
		UserID:    userID,
		GroupRole: role,
	}
	if err := r.db.WithContext(ctx).Create(&groupUser).Error; err != nil {
		return e.NewInternalServerApiError("Error adding user to group", err)
	}
	return nil
}

func (r *groupRepository) RemoveUserFromGroup(ctx context.Context, groupID int, userID int) e.ApiError {
	if err := r.db.WithContext(ctx).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Delete(&model.GroupXUsers{}).Error; err != nil {
		return e.NewInternalServerApiError("Error removing user from group", err)
	}
	return nil
}

func (r *groupRepository) UserExistsInGroup(ctx context.Context, groupID int, userID int) (bool, e.ApiError) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.GroupXUsers{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Count(&count).Error; err != nil {
		return false, e.NewInternalServerApiError("Error checking user in group", err)
	}
	return count > 0, nil
}

func generateGroupCode() string {
	rand.Seed(time.Now().UnixNano())
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 8)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}