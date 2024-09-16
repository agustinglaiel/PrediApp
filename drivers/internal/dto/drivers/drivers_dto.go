package dto

import "time"

type CreateDriverDTO struct {
    BroadcastName  string `json:"broadcast_name" binding:"required"`
    CountryCode    string `json:"country_code" binding:"required"`
    DriverNumber   int    `json:"driver_number" binding:"required"`
    FirstName      string `json:"first_name" binding:"required"`
    LastName       string `json:"last_name" binding:"required"`
    FullName       string `json:"full_name" binding:"required"`
    NameAcronym    string `json:"name_acronym" binding:"required"`
    TeamName       string `json:"team_name" binding:"required"`
}


type UpdateDriverDTO struct {
    BroadcastName  string `json:"broadcast_name,omitempty"`
    CountryCode    string `json:"country_code,omitempty"`
    DriverNumber   int    `json:"driver_number,omitempty"`
    FirstName      string `json:"first_name,omitempty"`
    LastName       string `json:"last_name,omitempty"`
    FullName       string `json:"full_name,omitempty"`
    NameAcronym    string `json:"name_acronym,omitempty"`
    TeamName       string `json:"team_name,omitempty"`
}

type ResponseDriverDTO struct {
    ID             uint   `json:"driver_id"`
    BroadcastName  string `json:"broadcast_name"`
    CountryCode    string `json:"country_code"`
    DriverNumber   int    `json:"driver_number"`
    FirstName      string `json:"first_name"`
    LastName       string `json:"last_name"`
    FullName       string `json:"full_name"`
    NameAcronym    string `json:"name_acronym"`
    TeamName       string `json:"team_name"`
}

type DeleteDriverDTO struct {
    DriverID uint `json:"driver_id"`
}

type ResponseEventDTO struct {
	ID                uint      `json:"id"`
	SessionID         uint      `json:"session_id"`
	RaceResultID      uint      `json:"race_result_id,omitempty"`
	SprintRaceResultID uint     `json:"sprint_race_result_id,omitempty"`
	QualyResultID     uint      `json:"qualy_result_id,omitempty"`
	SprintQualyResultID uint    `json:"sprint_qualy_result_id,omitempty"`
	FP1ID             uint      `json:"fp1_id,omitempty"`
	FP2ID             uint      `json:"fp2_id,omitempty"`
	FP3ID             uint      `json:"fp3_id,omitempty"`
	Date              time.Time `json:"date"`
}