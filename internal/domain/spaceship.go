package domain

// custom type for spaceship status enum
type SpaceshipStatus uint

const (
	// since iota starts with 0, the first value reserved for undefined
	SpaceshipStatusUndefined SpaceshipStatus = iota
	SpaceshipStatusOperational
	SpaceshipStatusDamaged
)

// convert status to string value
func (s SpaceshipStatus) String() string {
	return [...]string{
		"Undefined",
		"Operational",
		"Damaged",
	}[s]
}

// TODO: redo enum
func SpaceshipStatusFromString(s string) SpaceshipStatus {
	switch s {
	case "operational":
		return SpaceshipStatusOperational
	case "damaged":
		return SpaceshipStatusDamaged
	default:
		return SpaceshipStatusUndefined
	}
}

// spaceship armament with qty
type SpaceshipArmament struct {
	ID    uint
	Title string
	Qty   uint
}

// main spaceship model
type Spaceship struct {
	ID        uint
	Name      string
	Class     string
	Armament  []SpaceshipArmament
	Crew      uint
	Image     string
	Value     float64
	Status    SpaceshipStatus
	CreatedAt int64
	UpdatedAt int64
}
