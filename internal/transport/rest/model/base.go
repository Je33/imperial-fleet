package model

type PostResponce struct {
	Success bool `json:"success"`
}

type SpaceshipsResponce struct {
	Data []SpaceshipShort `json:"data"`
}
