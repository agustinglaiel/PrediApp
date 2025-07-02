package repository

import (
	"context"
	"math/rand"
	"time"

	"prediapp.local/db/model"
	e "prediapp.local/groups/pkg/utils"

	"gorm.io/gorm"
)

type groupRepository struct {
	db *gorm.DB
}

type GroupRepository interface {
	CreateGroup(ctx context.Context, group *model.Group) e.ApiError
	GetGroupByID(ctx context.Context, id int) (*model.Group, e.ApiError)
	GetGroupsByUserId(ctx context.Context, userID int) ([]model.Group, e.ApiError)
	GetGroupByGroupName(ctx context.Context, groupName string) (*model.Group, e.ApiError)
	GetGroupByCode(ctx context.Context, code string) (*model.Group, e.ApiError)
	GetGroups(ctx context.Context) ([]*model.Group, e.ApiError)
	DeleteGroupByID(ctx context.Context, id int) e.ApiError
	DeleteGroupByGroupName(ctx context.Context, groupName string) e.ApiError

	// Manejo de `groups_users`
	AddUserToGroup(ctx context.Context, groupID int, userID int, role string) e.ApiError
	RemoveUserFromGroup(ctx context.Context, groupID int, userID int) e.ApiError
	UserExistsInGroup(ctx context.Context, groupID int, userID int) (bool, e.ApiError)
	GetUserRoleInGroup(ctx context.Context, groupID int, userID int) (string, e.ApiError)
	UpdateUserRoleInGroup(ctx context.Context, groupID int, userID int, newRole string) e.ApiError
	GetJoinRequests(ctx context.Context, groupID int) ([]string, e.ApiError)
}

func NewGroupRepository(db *gorm.DB) GroupRepository {
	return &groupRepository{db: db}
}

func (r *groupRepository) CreateGroup(ctx context.Context, group *model.Group) e.ApiError {
	for {
		group.GroupCode = generateGroupCode() // Generar código aleatorio
		existingGroup, _ := r.GetGroupByCode(ctx, group.GroupCode)
		if existingGroup == nil {
			break
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
			return nil, e.NewNotFoundApiError("Grupo no encontrado.")
		}
		return nil, e.NewInternalServerApiError("Error finding group by ID", err)
	}
	return &group, nil
}

func (r *groupRepository) GetGroupsByUserId(ctx context.Context, userID int) ([]model.Group, e.ApiError) {
	var groups []model.Group
	err := r.db.WithContext(ctx).
		Preload("GroupUsers").
		Joins("JOIN group_x_users ON group_x_users.group_id = groups.id").
		Where("group_x_users.user_id = ?", userID).
		Find(&groups).Error // ⬅️ usamos Find (slice) en lugar de First
	if err != nil {
		return nil, e.NewInternalServerApiError("Error finding groups by user ID", err)
	}
	if len(groups) == 0 {
		return nil, e.NewNotFoundApiError("No groups found for user")
	}
	return groups, nil
}

func (r *groupRepository) GetGroupByGroupName(ctx context.Context, groupName string) (*model.Group, e.ApiError) {
	var group model.Group
	if err := r.db.WithContext(ctx).Where("group_name = ?", groupName).First(&group).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, e.NewNotFoundApiError("Grupo no encontrado.")
		}
		return nil, e.NewInternalServerApiError("error finding group by group name", err)
	}
	return &group, nil
}

func (r *groupRepository) GetGroupByCode(ctx context.Context, code string) (*model.Group, e.ApiError) {
	var group model.Group
	if err := r.db.WithContext(ctx).
		Preload("GroupUsers").
		Where("group_code = ?", code).
		First(&group).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, e.NewNotFoundApiError("Grupo no encontrado.")
		}
		return nil, e.NewInternalServerApiError("Error finding group by code", err)
	}
	return &group, nil
}

func (r *groupRepository) GetGroups(ctx context.Context) ([]*model.Group, e.ApiError) {
	var groups []*model.Group
	if err := r.db.WithContext(ctx).
		Preload("GroupUsers.User").
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

// Obtener el rol actual de un usuario en un grupo
func (r *groupRepository) GetUserRoleInGroup(ctx context.Context, groupID int, userID int) (string, e.ApiError) {
	var role string
	err := r.db.WithContext(ctx).
		Table("group_x_users").
		Select("group_role").
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Scan(&role).Error

	if err != nil {
		return "", e.NewInternalServerApiError("Error getting user role in group", err)
	}

	return role, nil
}

// Actualizar el rol de un usuario en un grupo
func (r *groupRepository) UpdateUserRoleInGroup(ctx context.Context, groupID int, userID int, newRole string) e.ApiError {
	err := r.db.WithContext(ctx).
		Model(&model.GroupXUsers{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Update("group_role", newRole).Error

	if err != nil {
		return e.NewInternalServerApiError("Error updating user role in group", err)
	}

	return nil
}

func (r *groupRepository) GetJoinRequests(ctx context.Context, groupID int) ([]string, e.ApiError) {
	var usernames []string
	err := r.db.WithContext(ctx).
		Table("group_x_users").
		Joins("JOIN users ON users.id = group_x_users.user_id").
		Select("users.username").
		Where("group_x_users.group_id = ? AND group_x_users.group_role = ?", groupID, "invited").
		Scan(&usernames).Error

	if err != nil {
		return nil, e.NewInternalServerApiError("Error getting join requests", err)
	}

	return usernames, nil
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
