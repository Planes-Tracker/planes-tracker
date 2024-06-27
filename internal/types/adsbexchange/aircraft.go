package adsbexchange

type navModesT uint8

const (
	NavModeAutopilot navModesT = 1 << iota
	NavModeVNAV
	NavModeAltHold
	NavModeApproach
	NavModeLNAV
	NavModeTCAS
)

type emergencyT uint8

const (
	EmergencyNone emergencyT = iota
	EmergencyGeneral
	EmergencyLifeguard
	EmergencyMinFuel
	EmergencyNORDO
	EmergencyUnlawful
	EmergencyDowned
	EmergencyReserved
)

type navAltitudeSourceT uint8

const (
	NavAltInvalid navAltitudeSourceT = iota
	NavAltUnknown
	NavAltAircraft
	NavAltMCP
	NavAltFMS
)

type airgroundT uint8

const (
	AGInvalid airgroundT = iota
	AGGround
	AGAirborne
	AGUncertain
)

type silTypeT uint8

const (
	SILInvalid silTypeT = iota
	SILUnknown
	SILPerSample
	SILPerHour
)

type addrtypeT uint8

const (
	AddrADSBICAO addrtypeT = iota
	AddrADSBICAONT
	AddrADSRICAO
	AddrTISBICAO
	AddrJAERO
	AddrMLAT
	AddrOther
	AddrModeS
	AddrADSBOther
	AddrADSROther
	AddrTISBTrackFile
	AddrTISBOther
	AddrModeA
	AddrUnknown
)

// https://github.com/ADSBexchange/readsb/blob/eea98b13b453779b9fe75b6783f4e6b129f2dff0/aircraft.h#L55-L184
type BinaryAircraft struct {
	Hex                    uint32
	Seen                   int32
	Lon                    int32
	Lat                    int32
	BaroRate               int16
	GeomRate               int16
	BaroAlt                int16
	GeomAlt                int16
	NavAltitudeMCP         uint16
	NavAltitudeFMS         uint16
	NavQNH                 int16
	NavHeading             int16
	Squawk                 uint16
	GS                     int16
	Mach                   int16
	Roll                   int16
	Track                  int16
	TrackRate              int16
	MagHeading             int16
	TrueHeading            int16
	WindDirection          int16
	WindSpeed              int16
	OAT                    int16
	TAT                    int16
	TAS                    uint16
	IAS                    uint16
	PosRC                  uint16
	Messages               uint16
	Category               uint8
	PosNIC                 uint8
	NavModes               navModesT
	Emergency              emergencyT         `bitfield:"4"`
	AddrType               addrtypeT          `bitfield:"4"`
	Airground              airgroundT         `bitfield:"4"`
	NavAltitudeSource      navAltitudeSourceT `bitfield:"4"`
	SilType                silTypeT           `bitfield:"4"`
	ADSBVersion            uint8              `bitfield:"4"`
	ADSRVersion            uint8              `bitfield:"4"`
	TISBVersion            uint8              `bitfield:"4"`
	NACP                   uint8              `bitfield:"4"`
	NACV                   uint8              `bitfield:"4"`
	SIL                    uint8              `bitfield:"2"`
	GVA                    uint8              `bitfield:"2"`
	SDA                    uint8              `bitfield:"2"`
	NIC_A                  uint8              `bitfield:"1"`
	NIC_C                  uint8              `bitfield:"1"`
	NICBaro                uint8              `bitfield:"1"`
	Alert                  uint8              `bitfield:"1"`
	SPI                    uint8              `bitfield:"1"`
	CallsignValid          uint8              `bitfield:"1"`
	BaroAltValid           uint8              `bitfield:"1"`
	GeomAltValid           uint8              `bitfield:"1"`
	PositionValid          uint8              `bitfield:"1"`
	GSValid                uint8              `bitfield:"1"`
	IASValid               uint8              `bitfield:"1"`
	TASValid               uint8              `bitfield:"1"`
	MachValid              uint8              `bitfield:"1"`
	TrackValid             uint8              `bitfield:"1"`
	TrackRateValid         uint8              `bitfield:"1"`
	RollValid              uint8              `bitfield:"1"`
	MagHeadingValid        uint8              `bitfield:"1"`
	TrueHeadingValid       uint8              `bitfield:"1"`
	BaroRateValid          uint8              `bitfield:"1"`
	GeomRateValid          uint8              `bitfield:"1"`
	NIC_AValid             uint8              `bitfield:"1"`
	NIC_CValid             uint8              `bitfield:"1"`
	NICBaroValid           uint8              `bitfield:"1"`
	NACPValid              uint8              `bitfield:"1"`
	NACVValid              uint8              `bitfield:"1"`
	SILValid               uint8              `bitfield:"1"`
	GVAValid               uint8              `bitfield:"1"`
	SDAValid               uint8              `bitfield:"1"`
	SquawkValid            uint8              `bitfield:"1"`
	EmergencyValid         uint8              `bitfield:"1"`
	SPIValid               uint8              `bitfield:"1"`
	NavQNHValid            uint8              `bitfield:"1"`
	NavAltitudeMCPValid    uint8              `bitfield:"1"`
	NavAltitudeFMSValid    uint8              `bitfield:"1"`
	NavAltitudeSrcValid    uint8              `bitfield:"1"`
	NavHeadingValid        uint8              `bitfield:"1"`
	NavModesValid          uint8              `bitfield:"1"`
	AlertValid             uint8              `bitfield:"1"`
	WindValid              uint8              `bitfield:"1"`
	TempValid              uint8              `bitfield:"1"`
	Unused1                uint8              `bitfield:"1"`
	Unused2                uint8              `bitfield:"1"`
	Callsign               [8]byte
	DBFlags                uint16
	TypeCode               [4]byte
	Registration           [12]byte
	ReceiverCount          uint8
	Signal                 uint8
	ExtraFlags             uint8
	Reserved               uint8
	SeenPost               int32
}
