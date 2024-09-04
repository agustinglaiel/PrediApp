package model

import (
	"time"
)

type Event struct {
	ID                  int        `gorm:"primary_key;auto_increment" json:"id"`
	SessionID           int        `gorm:"not null" json:"session_id"`
	RaceResultID        int       `json:"race_result_id"`
	SprintRaceResultID  int       `json:"sprint_race_result_id"`
	QualyResultID       int       `json:"qualy_result_id"`
	SprintQualyResultID int       `json:"sprint_qualy_result_id"`
	FP1ID               int       `json:"fp1_id"`
	FP2ID               int       `json:"fp2_id"`
	FP3ID               int       `json:"fp3_id"`
	Date                time.Time  `gorm:"not null" json:"date"`
}
