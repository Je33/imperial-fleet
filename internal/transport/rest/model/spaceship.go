package model

import "github.com/Je33/imperial_fleet/internal/domain"

type SpaceshipShort struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type SpaceshipFull struct {
	ID       uint                       `json:"id"`
	Name     string                     `json:"name"`
	Class    string                     `json:"class"`
	Armament []domain.SpaceshipArmament `json:"armament"`
	Crew     uint                       `json:"crew"`
	Image    string                     `json:"image"`
	Value    float64                    `json:"value"`
	Status   string                     `json:"status"`
}


