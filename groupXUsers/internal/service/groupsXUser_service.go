package service

import (
	"context"
	"groupXUsers/internal/dto"
	"groupXUsers/internal/model"
	"groupXUsers/internal/repository"
	e "groupXUsers/pkg/utils"
)


type groupXUserService struct {
	groupXUserRepo repository.GroupXUserRepository
}

type GroupXUserServiceInterface interface {
	CreateGroupXUser(ctx context.Context, groupXUser dto.GroupXUserCreateDTO) (dto.GroupXUserCreateDTO, e.ApiError)
	GetGroupXUserByGroupID(ctx context.Context, groupId int) ([]dto.GroupXUserResponseDTO, e.ApiError)
	GetGroupXUserByGroupName(ctx context.Context, groupName string) ([]dto.GroupXUserResponseDTO, e.ApiError)
	GetGroupXUsers(ctx context.Context) ([]dto.GroupXUserResponseDTO, e.ApiError)
	RemoveUserFromGroup(ctx context.Context, userId int, groupId int) e.ApiError
	RemoveUserFromGroupByGroupName(ctx context.Context, userId int, groupName string) e.ApiError
}

func NewGroupXUserService(groupXUserRepo repository.GroupXUserRepository) GroupXUserServiceInterface {
	return &groupXUserService{groupXUserRepo: groupXUserRepo}
}

func (s *groupXUserService) CreateGroupXUser(ctx context.Context, groupXUser dto.GroupXUserCreateDTO) (dto.GroupXUserCreateDTO, e.ApiError) {
	// Check if the groupXUser already exists
	if _, err := s.groupXUserRepo.GetGroupXUserByGroupID(ctx, groupXUser.GroupId); err == nil {
		return dto.GroupXUserCreateDTO{}, e.NewBadRequestApiError("groupXUser already exists")
	}

	newGroupXUser := &model.GroupXUser{
		UserId: groupXUser.UserId,
		GroupId: groupXUser.GroupId,
	}

	if err := s.groupXUserRepo.CreateGroupXUser(ctx, newGroupXUser); err != nil {
		return dto.GroupXUserCreateDTO{}, e.NewInternalServerApiError("error creating groupXUser", err)
	}

	response := dto.GroupXUserCreateDTO{
		UserId: newGroupXUser.UserId,
		GroupId: newGroupXUser.GroupId,
	}

	return response, nil
}

