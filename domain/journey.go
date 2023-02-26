package domain

import (
	"time"

	"github.com/google/uuid"
)

type Journey struct {
	JourneyId              int64
	TripId                 uuid.UUID
	JourneyStartDatetime   time.Time
	JourneyStartDate       time.Time
	JourneyStartTime       time.Time
	JourneyStartLon        uint64
	JourneyStartLat        uint64
	JourneyStartInsee      int64
	JourneyStartPostalcode string
	JourneyStartDepartment string
	JourneyStartTown       string
	JourneyStartTowngroup  string
	JourneyStartCountry    string
	JourneyEndDatetime     time.Time
	JourneyEndDate         time.Time
	JourneyEndTime         time.Time
	JourneyEndLon          uint64
	JourneyEndLat          uint64
	JourneyEndInsee        int64
	JourneyEndPostalcode   string
	JourneyEndDepartment   string
	JourneyEndTown         string
	JourneyEndTowngroup    string
	JourneyEndCountry      string
	PassengerSeats         int16
	OperatorClass          string
	JourneyDistance        int64
	JourneyDuration        int64
	HasIncentive           bool
}
