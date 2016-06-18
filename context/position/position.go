package position

type Position int8

const (
	UNSUPPORTED Position = iota
	TOP
	JUNGLE
	MID
	ADCARRY
	SUPPORT
)

func (p Position) String() string {
	switch p {
	case TOP:
		return "TOP_LANE"
	case JUNGLE:
		return "JUNGLE"
	case MID:
		return "MID_LANE"
	case ADCARRY:
		return "AD_CARRY"
	case SUPPORT:
		return "SUPPORT"
	default:
		return "Unsupported position type"
	}
}

func FromString(s string) Position {
	switch s {
	case "TOP_LANE":
		return TOP
	case "JUNGLER":
		return JUNGLE
	case "MID_LANE":
		return MID
	case "AD_CARRY":
		return ADCARRY
	case "SUPPORT":
		return SUPPORT
	default:
		return UNSUPPORTED
	}
}
