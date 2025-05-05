package constants

type FieldScheduleStatusName string
type FieldScheduleStatus int

const (
	AvailableString FieldScheduleStatusName = "Available"
	BookedString    FieldScheduleStatusName = "Booked"

	Available FieldScheduleStatus = 100
	Booked    FieldScheduleStatus = 200
)

var mapFieldScheduleIntToString = map[FieldScheduleStatus]FieldScheduleStatusName{
	Available: AvailableString,
	Booked:    BookedString,
}

var mapFieldScheduleStringToInt = map[FieldScheduleStatusName]FieldScheduleStatus{
	AvailableString: Available,
	BookedString:    Booked,
}

func (f FieldScheduleStatus) GetStatusString() FieldScheduleStatusName {
	return mapFieldScheduleIntToString[f]
}

func (f FieldScheduleStatusName) GetStatusInt() FieldScheduleStatus {
	return mapFieldScheduleStringToInt[f]
}
