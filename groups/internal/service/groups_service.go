package service

import (
	"context"
	"groups/internal/dto"
	"groups/internal/model"
	"groups/internal/repository"
	e "groups/pkg/utils"
)


type groupService struct {
	groupRepo repository.GroupRepository
}

type GroupServiceInterface interface {
	CreateGroup(ctx context.Context, group dto.CreateGroupRequestDTO) (dto.GroupResponseDTO, e.ApiError)
	GetGroupByID(ctx context.Context, id int) (dto.GroupResponseDTO, e.ApiError)
	GetGroups(ctx context.Context) ([]dto.GroupListResponseDTO, e.ApiError)
	DeleteGroupByID(ctx context.Context, id int) e.ApiError
	JoinGroup(ctx context.Context, request dto.RequestJoinGroupDTO) e.ApiError
	ManageGroupInvitation(ctx context.Context, request dto.ManageGroupInvitationDTO) e.ApiError
}

func NewGroupService(groupRepo repository.GroupRepository) GroupServiceInterface {
	return &groupService{groupRepo: groupRepo}
}

func (s *groupService) CreateGroup(ctx context.Context, request dto.CreateGroupRequestDTO) (dto.GroupResponseDTO, e.ApiError) {
	// Verificar si el nombre del grupo ya existe
	if _, err := s.groupRepo.GetGroupByGroupName(ctx, request.GroupName); err == nil {
		return dto.GroupResponseDTO{}, e.NewBadRequestApiError("Group name already exists")
	}

	newGroup := &model.Group{
		GroupName:   request.GroupName,
		Description: request.Description,
	}

	// Crear el grupo en la base de datos
	if err := s.groupRepo.CreateGroup(ctx, newGroup); err != nil {
		return dto.GroupResponseDTO{}, e.NewInternalServerApiError("Error creating group", err)
	}

	// Asociar al usuario creador en `groups_users` con el rol de `creator`
	if err := s.groupRepo.AddUserToGroup(ctx, newGroup.ID, request.UserID, "creator"); err != nil {
		return dto.GroupResponseDTO{}, e.NewInternalServerApiError("Error adding creator to group", err)
	}

	// Construir la respuesta con el código único
	response := dto.GroupResponseDTO{
		ID:          newGroup.ID,
		GroupName:   newGroup.GroupName,
		Description: newGroup.Description,
		GroupCode:   newGroup.GroupCode,
		Users:       []dto.GroupUserResponseDTO{}, // No necesitamos cargar los usuarios aquí
	}

	return response, nil
}

func (s *groupService) GetGroupByID(ctx context.Context, id int) (dto.GroupResponseDTO, e.ApiError) {
	group, err := s.groupRepo.GetGroupByID(ctx, id)
	if err != nil {
		return dto.GroupResponseDTO{}, err
	}

	// Mapear usuarios del grupo
	var users []dto.GroupUserResponseDTO
	for _, groupUser := range group.GroupUsers { 
		users = append(users, dto.GroupUserResponseDTO{
			UserID: groupUser.UserID,
			Role:   groupUser.GroupRole,
			Score:  nil, // Se traerá desde el microservicio de `users` en otra función
		})
	}

	response := dto.GroupResponseDTO{
		ID:          group.ID,
		GroupName:   group.GroupName,
		Description: group.Description,
		GroupCode:   group.GroupCode,
		Users:       users,
	}

	return response, nil
}


func (s *groupService) GetGroups(ctx context.Context) ([]dto.GroupListResponseDTO, e.ApiError) {
	groups, err := s.groupRepo.GetGroups(ctx)
	if err != nil {
		return nil, err
	}

	var response []dto.GroupListResponseDTO
	for _, group := range groups {
		response = append(response, dto.GroupListResponseDTO{
			ID:        group.ID,
			GroupName: group.GroupName,
			GroupCode: group.GroupCode,
		})
	}

	return response, nil
}

func (s *groupService) DeleteGroupByID(ctx context.Context, id int) e.ApiError {
	_, err := s.groupRepo.GetGroupByID(ctx, id)
	if err != nil {
		return e.NewBadRequestApiError("Group not found")
	}

	if err := s.groupRepo.DeleteGroupByID(ctx, id); err != nil {
		return e.NewInternalServerApiError("Error deleting group", err)
	}
	return nil
}

func (s *groupService) DeleteGroupByGroupName(ctx context.Context, groupName string) e.ApiError {
	_, err := s.groupRepo.GetGroupByGroupName(ctx, groupName)
	if err != nil {
		return e.NewBadRequestApiError("group not found")
	}

	if err := s.groupRepo.DeleteGroupByGroupName(ctx, groupName); err != nil {
		return e.NewInternalServerApiError("error deleting group", err)
	}
	return nil
}

func (s *groupService) JoinGroup(ctx context.Context, request dto.RequestJoinGroupDTO) e.ApiError {
	group, err := s.groupRepo.GetGroupByCode(ctx, request.GroupCode)
	if err != nil {
		return err
	}

	// Verificar si el usuario ya está en el grupo
	exists, err := s.groupRepo.UserExistsInGroup(ctx, group.ID, request.UserID)
	if err != nil {
		return err
	}
	if exists {
		return e.NewBadRequestApiError("User is already in the group")
	}

	// Agregar al usuario como `invited`
	if err := s.groupRepo.AddUserToGroup(ctx, group.ID, request.UserID, "invited"); err != nil {
		return e.NewInternalServerApiError("Error adding user to group", err)
	}

	return nil
}

func (s *groupService) ManageGroupInvitation(ctx context.Context, request dto.ManageGroupInvitationDTO) e.ApiError {
	group, err := s.groupRepo.GetGroupByID(ctx, request.GroupID)
	if err != nil {
		return err
	}

	// Verificar si el usuario que realiza la acción es el creador del grupo
	isCreator, err := s.groupRepo.UserExistsInGroup(ctx, group.ID, request.UserID)
	if err != nil {
		return err
	}
	if !isCreator {
		return e.NewForbiddenApiError("Only the group creator can manage invitations")
	}

	// Aceptar o rechazar la invitación
	if request.Action == "accept" {
		if err := s.groupRepo.AddUserToGroup(ctx, group.ID, request.UserID, "member"); err != nil {
			return e.NewInternalServerApiError("Error accepting user into group", err)
		}
	} else if request.Action == "reject" {
		if err := s.groupRepo.RemoveUserFromGroup(ctx, group.ID, request.UserID); err != nil {
			return e.NewInternalServerApiError("Error rejecting user from group", err)
		}
	} else {
		return e.NewBadRequestApiError("Invalid action")
	}

	return nil
}