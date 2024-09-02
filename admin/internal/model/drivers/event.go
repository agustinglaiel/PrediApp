package model

import "time"

type Event struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
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