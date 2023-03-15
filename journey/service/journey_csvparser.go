package service

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
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

type job struct {
	line []string
	lineNumber int
}

func (p *journeyCsvParser) Parse(reader io.Reader, journeyChan chan<- *domain.Journey, errorChan chan<- string) {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = ';'

	p.logger.Debug("Reading headers")
	if _, err := csvReader.Read(); err != nil {
		
		close(journeyChan)
		if err.Error() == "EOF" {
			p.logger.Info("End of the file")
			close(errorChan)
			return
		}

		p.logger.Errorw("Error reading headers",
			"error", err,
		)
		errorChan <- err.Error()
		close(errorChan)
		return
	}

	numWorkers := 1
	jobs := make(chan *job, numWorkers)

	var workerGroup sync.WaitGroup

	worker := func(jobs <-chan *job, results chan<- *domain.Journey, errorChan chan<- string) {
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

					res, err := p.parseJourney(job.line, job.lineNumber);
					if  err != nil{
						errorChan <- err.Error()
					} else {
						results <- res
					}
			}
		}
	}

	p.logger.Debug("Workers Initializing")
	for w := 0; w < numWorkers; w++ {
		workerGroup.Add(1)
		go func() {
			defer workerGroup.Done()
			worker(jobs, journeyChan, errorChan)
		}()
	}

	lineNumber := 0
	go func() {
		for {
			lineNumber += 1
			line, err := csvReader.Read()
			if err == io.EOF {
				p.logger.Debug("End of file reached")
				break
			}

			if err != nil {
				p.logger.Error("Error reading csv file:", err.Error())
				errorChan <- err.Error()
				break
			}
			jobs <- &job{
				line: line,
				lineNumber: lineNumber,
			}
		}
		close(jobs)
	}()

	go func() {
		workerGroup.Wait()
		p.logger.Debug("Closing channels")
		close(journeyChan)
		close(errorChan)
	}()

}

func (p *journeyCsvParser) parseJourney(r []string, lineNumber int) (journey *domain.Journey, err error) {
	defer func(){
		if(recover() != nil){
			err = fmt.Errorf("problem while parsing a journey: line %d in wrong format", lineNumber)
			p.logger.Errorw(err.Error(),
				"line", strings.Join(r, ","),
			)
		}
	}()

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

	journey = &domain.Journey{
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

	return journey, nil
}
