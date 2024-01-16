package model

type SpaceshipArmament struct {
	Title string `json:"title"`
	Qty   uint   `json:"qty,string"`
}

type SpaceshipShort struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type SpaceshipFull struct {
	ID       uint                `json:"id"`
	Name     string              `json:"name"`
	Class    string              `json:"class"`
	Armament []SpaceshipArmament `json:"armament"`
	Crew     uint                `json:"crew"`
	Image    string              `json:"image"`
	Value    float64             `json:"value"`
	Status   string              `json:"status"`
}
