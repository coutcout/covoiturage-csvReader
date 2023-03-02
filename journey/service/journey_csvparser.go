package service

import (
	"encoding/csv"
	"io"
	"strconv"
	"sync"
	"time"

	"me/coutcout/covoiturage/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type journeyCsvParser struct {
	logger *zap.SugaredLogger
}

func NewJourneyCsvParser(logger *zap.SugaredLogger) domain.JourneyParser {
	return &journeyCsvParser{logger}
}

func (p *journeyCsvParser) Parse(reader io.Reader, journeyChan chan<- *domain.Journey) error {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = ';'

	p.logger.Debug("Reading headers")
	if _, err := csvReader.Read(); err != nil {
		p.logger.Error("Error reading headers",
			"error", err,
		)
		close(journeyChan)
		if err.Error() == "EOF" {
			return nil
		}

		return err
	}

	numWorkers := 10
	jobs := make(chan []string, numWorkers)

	var workerGroup sync.WaitGroup

	worker := func(jobs <-chan []string, results chan<- *domain.Journey) {
		p.logger.Debug("Worker started")
		for {
			select {
			case job, ok := <-jobs:
				if !ok {
					p.logger.Debug("Worker ended")
					return
				}
				p.logger.Debugw("Line received",
					"csvLine", job,
				)
				results <- p.parseJourney(job)
			}
		}
	}

	p.logger.Debug("Workers Initializing")
	for w := 0; w < numWorkers; w++ {
		workerGroup.Add(1)
		go func() {
			defer workerGroup.Done()
			worker(jobs, journeyChan)
		}()
	}

	go func() {
		for {
			line, err := csvReader.Read()
			if err == io.EOF {
				p.logger.Debug("End of file reached")
				break
			}

			if err != nil {
				p.logger.Error("Error reading csv file:", err.Error())
				break
			}
			jobs <- line
		}
		close(jobs)
	}()

	go func() {
		workerGroup.Wait()
		p.logger.Debug("Closing channel")
		close(journeyChan)
	}()

	return nil
}

func (p *journeyCsvParser) parseJourney(r []string) *domain.Journey {
	journeyId, _ := strconv.ParseInt(r[0], 10, 64)
	tripId, _ := uuid.Parse(r[1])
	startDateTime, _ := time.Parse("2006-01-02T15:04:05-07:00", r[2])
	startDate, _ := time.Parse(time.DateOnly, r[3])
	startTime, _ := time.Parse(time.TimeOnly, r[4])
	startLon, _ := strconv.ParseUint(r[5], 10, 64)
	startLat, _ := strconv.ParseUint(r[6], 10, 64)
	startInsee, _ := strconv.ParseInt(r[7], 10, 64)

	endDateTime, _ := time.Parse("2006-01-02T15:04:05-07:00", r[13])
	endDate, _ := time.Parse(time.DateOnly, r[14])
	endTime, _ := time.Parse(time.TimeOnly, r[15])
	endLon, _ := strconv.ParseUint(r[16], 10, 64)
	endLat, _ := strconv.ParseUint(r[17], 10, 64)
	endInsee, _ := strconv.ParseInt(r[18], 10, 64)
	passagerSeats, _ := strconv.ParseInt(r[24], 10, 16)
	distance, _ := strconv.ParseInt(r[26], 10, 64)
	duration, _ := strconv.ParseInt(r[27], 10, 64)
	hasIncentive := r[28] == "OUI"

	journey := &domain.Journey{
		JourneyId:              journeyId,
		TripId:                 tripId,
		JourneyStartDatetime:   startDateTime,
		JourneyStartDate:       startDate,
		JourneyStartTime:       startTime,
		JourneyStartLon:        startLon,
		JourneyStartLat:        startLat,
		JourneyStartInsee:      startInsee,
		JourneyStartPostalcode: r[8],
		JourneyStartDepartment: r[9],
		JourneyStartTown:       r[10],
		JourneyStartTowngroup:  r[11],
		JourneyStartCountry:    r[12],
		JourneyEndDatetime:     endDateTime,
		JourneyEndDate:         endDate,
		JourneyEndTime:         endTime,
		JourneyEndLon:          endLon,
		JourneyEndLat:          endLat,
		JourneyEndInsee:        endInsee,
		JourneyEndPostalcode:   r[19],
		JourneyEndDepartment:   r[20],
		JourneyEndTown:         r[21],
		JourneyEndTowngroup:    r[22],
		JourneyEndCountry:      r[23],
		PassengerSeats:         int16(passagerSeats),
		OperatorClass:          r[25],
		JourneyDistance:        distance,
		JourneyDuration:        duration,
		HasIncentive:           hasIncentive,
	}
	return journey
}
