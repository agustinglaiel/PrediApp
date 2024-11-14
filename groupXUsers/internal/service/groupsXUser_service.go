package service

import (
	"context"
	"groupXUsers/internal/client"
	"groupXUsers/internal/dto"
	"groupXUsers/internal/model"
	"groupXUsers/internal/repository"
	e "groupXUsers/pkg/utils"
)


type groupXUserService struct {
	groupXUserRepo repository.GroupXUserRepository
	client         *client.HttpClient
}

type GroupXUserServiceInterface interface {
	CreateGroupXUser(ctx context.Context, groupXUser dto.GroupXUserCreateDTO) (dto.GroupXUserCreateDTO, e.ApiError)
	GetGroupXUserByGroupID(ctx context.Context, groupId int) ([]dto.GroupXUserResponseDTO, e.ApiError)
	GetGroupXUserByGroupName(ctx context.Context, groupName string) ([]dto.GroupXUserResponseDTO, e.ApiError)
	GetGroupXUsers(ctx context.Context) ([]dto.GroupXUserResponseDTO, e.ApiError)
	RemoveUserFromGroup(ctx context.Context, userId int, groupId int) e.ApiError
}

func NewGroupXUserService(groupXUserRepo repository.GroupXUserRepository, client *client.HttpClient) GroupXUserServiceInterface {
    return &groupXUserService{
        groupXUserRepo: groupXUserRepo,
		client: client,
    }
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

func (s *groupXUserService) GetGroupXUserByGroupID(ctx context.Context, groupId int) ([]dto.GroupXUserResponseDTO, e.ApiError) {
	// Check if the group exists using groupRepo function GetGroupByID
	groupId, err := s.client.GetGroupByID(groupId)
	if err != nil {
		return nil, e.NewNotFoundApiError("group not found")
	}
	
	groupXUsers, err := s.groupXUserRepo.GetGroupXUserByGroupID(ctx, groupId)
	if err != nil {
		return []dto.GroupXUserResponseDTO{}, e.NewBadRequestApiError("error getting groupXUser by groupID")
	}

	var response []dto.GroupXUserResponseDTO
	for _, groupXUser := range groupXUsers {
		response = append(response, dto.GroupXUserResponseDTO{
			UserId: groupXUser.UserId,
			GroupId: groupXUser.GroupId,
		})
	}

	return response, nil
}

func (s *groupXUserService) GetGroupXUserByGroupName(ctx context.Context, groupName string) ([]dto.GroupXUserResponseDTO, e.ApiError) {
	// Check if the group exists using groupRepo function GetGroupByGroupName
	groupName, err := s.client.GetGroupByGroupName(groupName)
	if err != nil {
		return nil, e.NewNotFoundApiError("group not found")
	}
	
	groupXUsers, err := s.groupXUserRepo.GetGroupXUserByGroupName(ctx, groupName)
	if err != nil {
		return []dto.GroupXUserResponseDTO{}, e.NewBadRequestApiError("error getting groupXUser by groupName")
	}

	var response []dto.GroupXUserResponseDTO
	for _, groupXUser := range groupXUsers {
		response = append(response, dto.GroupXUserResponseDTO{
			UserId: groupXUser.UserId,
			GroupId: groupXUser.GroupId,
		})
	}

	return response, nil
}

func (s *groupXUserService) GetGroupXUsers(ctx context.Context) ([]dto.GroupXUserResponseDTO, e.ApiError) {
	groupXUsers, err := s.groupXUserRepo.GetGroupXUsers(ctx)
	if err != nil {
		return nil, err
	}

	var response []dto.GroupXUserResponseDTO
	for _, groupXUser := range groupXUsers {
		response = append(response, dto.GroupXUserResponseDTO{
			UserId: groupXUser.UserId,
			GroupId: groupXUser.GroupId,
		})
	}

	return response, nil
}

func (s *groupXUserService) RemoveUserFromGroup(ctx context.Context, userId int, groupId int) e.ApiError {
	// Check if the groupXUser exists
	_, err := s.groupXUserRepo.GetGroupXUserByGroupIDAndUserID(ctx, groupId, userId)
	if err != nil {
		return e.NewNotFoundApiError("groupXUser not found")
	}

	if err := s.groupXUserRepo.RemoveUserFromGroup(ctx, userId, groupId); err != nil {
		return e.NewInternalServerApiError("error removing user from group", err)
	}

	return nil
}
