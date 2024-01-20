//
// SPDX-License-Identifier: Apache-2.0
//
// HybridClient deals with Redis for MetaData and Influx for Core-data

package hybrid

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/edgexfoundry/edgex-go/internal/pkg/db"
	"github.com/edgexfoundry/edgex-go/internal/pkg/db/influx"
	"github.com/edgexfoundry/edgex-go/internal/pkg/infrastructure/redis"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/errors"
	model "github.com/edgexfoundry/go-mod-core-contracts/v3/models"
	"github.com/google/uuid"
)

type HybridClient struct {
	redisClient   *redis.Client
	influxClient  *influx.Client
	loggingClient logger.LoggingClient
}

func NewHybridClient(config db.Configuration, logger logger.LoggingClient) (*HybridClient, errors.EdgeX) {
	//Create Influx client
	influxClient, err := influx.NewClient(config, logger) // create the client
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindDatabaseError, "influx client creation failed", err)
	}
	err = influxClient.CreateBucket()
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindDatabaseError, "bucket creation failed", err)
	}

	//Check the environment variables for redis client
	dbhost := os.Getenv("STAGEGATE_DATABASE_HOST")
	if dbhost != "" {
		config.Host = dbhost
	}
	dbport := os.Getenv("STAGEGATE_DATABASE_PORT")
	if dbport != "" {
		config.Port, _ = strconv.Atoi(dbport)
	}
	//Create Redis client
	redis, err := redis.NewClient(config, logger)
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindDatabaseError, "redis client creation failed", err)
	}

	//Create the Hybrid client
	hybridClient := &HybridClient{redisClient: redis, influxClient: influxClient, loggingClient: logger}
	return hybridClient, nil
}

// CloseSession
func (c *HybridClient) CloseSession() {
	c.influxClient.CloseSession()
	c.redisClient.CloseSession()
}

// AddEvent adds a new event
func (c *HybridClient) AddEvent(e model.Event) (model.Event, errors.EdgeX) {
	if e.Id != "" {
		_, err := uuid.Parse(e.Id)
		if err != nil {
			return model.Event{}, errors.NewCommonEdgeX(errors.KindInvalidId, "uuid parsing failed", err)
		}
	}
	return AddEvent(c, e)
}

// EventById gets an event by id
func (c *HybridClient) EventById(id string) (event model.Event, edgeXerr errors.EdgeX) {
	events1, err := AllEvents(c, 0, 0,
		fmt.Sprintf(`
		|>range(start: 0, stop: now())
		|>group()
		|>sort(columns: ["_time"])
		|>pivot(rowKey: ["_time","counter","devicename","_measurement","resourcename"], columnKey: ["_field"], valueColumn: "_value")
		|>filter(fn: (r) => r["eventid"] == "%v")`, id))
	if err != nil {
		return model.Event{}, err
	}
	if len(events1) > 0 {
		return events1[0], nil

	}
	return model.Event{}, nil
}

// DeleteEventById removes an event by id
func (c *HybridClient) DeleteEventById(id string) (edgeXerr errors.EdgeX) {
	//Get the event
	event, err := c.EventById(id)
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}
	//Create the delete string
	deletestring := fmt.Sprintf(`_measurement="%v" AND devicename="%v"`, event.ProfileName, event.DeviceName)

	err1 := c.influxClient.DeleteData(deletestring, time.Unix(0, int64(event.Origin)), time.Unix(0, int64(event.Origin)))
	if err1 != nil {
		return errors.NewCommonEdgeXWrapper(err1)
	}
	return nil
}
func (c *HybridClient) EventTotalCount() (uint32, errors.EdgeX) {
	//Influx count
	count, err := TotalCountInflux(c,
		`
		|>range(start: 0, stop: now())
		|>group()
		|>filter(fn: (r) => r._field == "eventid")
		|>distinct()
		|>count()`)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// EventCountByDeviceName returns the count of Event associated a specific Device from the database
func (c *HybridClient) EventCountByDeviceName(deviceName string) (uint32, errors.EdgeX) {
	//Influx count
	count, err := TotalCountInflux(c,
		fmt.Sprintf(`
		|>range(start: 0, stop: now())
		|>filter(fn: (r) => r["devicename"] == "%v")
		|>group()
		|>filter(fn: (r) => r._field == "eventid")
		|>distinct()
		|>count()`, deviceName))
	if err != nil {
		return 0, err
	}

	return count, nil
}

// LatestReadingByOffset returns a latest reading by offset
func (c *HybridClient) LatestReadingByOffset(offset uint32) (model.Reading, errors.EdgeX) {
	readings, err := allInfluxReadings(c,
		fmt.Sprintf(`
		|>range(start: 0, stop: now())
		|>filter(fn: (r) => r["_field"] == "eventtype" or r["_field"] == "readingid" or r["_field"] == "units" or r["_field"] == "readingorigin" or r["_field"] == "valuetype" or r["_field"] == "value")
		|>group()
		|>sort(columns: ["_time"])
		|>pivot(rowKey: ["_time","devicename","_measurement","resourcename"], columnKey: ["_field"], valueColumn: "_value")
		|>limit(n: %v, offset: %v)`, offset-1, offset))
	if err != nil {
		return nil, err
	}
	if len(readings) > 0 {
		return readings[0], nil
	}
	return nil, nil
}

// EventCountByTimeRange returns the count of Event by time range
func (c *HybridClient) EventCountByTimeRange(startTime int, endTime int) (uint32, errors.EdgeX) {
	//Influx count
	count, err := TotalCountInflux(c,
		fmt.Sprintf(`
			|>range(start:time(v:%v), stop: time(v:%v))
			|>group()
			|>filter(fn: (r) => r._field == "eventid")
			|>distinct()
			|>count()`, startTime, endTime))
	if err != nil {
		return 0, err
	}

	return count, nil
}

// AllEvents query events by offset and limit
func (c *HybridClient) AllEvents(offset int, limit int) ([]model.Event, errors.EdgeX) {
	events1, err := AllEvents(c, offset, limit,
		`
		|>range(start: 0, stop: now())
		|>group()
		|>sort(columns: ["_time"])
		|>pivot(rowKey: ["_time","counter","devicename","_measurement","resourcename"], columnKey: ["_field"], valueColumn: "_value")`)
	if err != nil {
		return nil, err
	}
	return events1, nil
}

// DeleteEventsByDeviceName deletes specific device's events and corresponding readings.  This function is implemented to starts up
// two goroutines to delete readings and events in the background to achieve better performance.
func (c *HybridClient) DeleteEventsByDeviceName(deviceName string) (edgeXerr errors.EdgeX) {
	err := c.influxClient.DeleteData(fmt.Sprintf(`devicename="%v"`, deviceName), time.Now().AddDate(-10, 0, 0) /*10 years back*/, time.Now())
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}
	return nil
}

// DeleteEventsByAge deletes events and their corresponding readings that are older than age.  This function is implemented to starts up
func (c *HybridClient) DeleteEventsByAge(age int64) (edgeXerr errors.EdgeX) {
	// note that the origin time is in  Epoch timestamp/nanoseconds format
	expireTimestamp := time.Now().UnixNano() - age
	err := c.influxClient.DeleteData("", time.Unix(0, expireTimestamp), time.Now())
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}
	return nil
}

// EventsByDeviceName query events by offset, limit and device name
func (c *HybridClient) EventsByDeviceName(offset int, limit int, name string) (events []model.Event, edgeXerr errors.EdgeX) {
	events1, err := AllEvents(c, offset, limit,
		fmt.Sprintf(`
		|>range(start: 0, stop: now())
		|>filter(fn: (r) => r["_measurement"] == "%v")
		|>group()
		|>sort(columns: ["_time"])
		|>pivot(rowKey: ["_time","counter","devicename","_measurement","resourcename"], columnKey: ["_field"], valueColumn: "_value")`, name))
	if err != nil {
		return nil, err
	}
	return events1, nil
}

// EventsByTimeRange query events by time range, offset, and limit
func (c *HybridClient) EventsByTimeRange(startTime int, endTime int, offset int, limit int) (events []model.Event, edgeXerr errors.EdgeX) {
	events1, err := AllEvents(c, offset, limit,
		fmt.Sprintf(`
			|> range(start:time(v:%v), stop: time(v:%v))
			|>group()
			|>sort(columns: ["_time"])
			|>pivot(rowKey: ["_time","counter","devicename","_measurement","resourcename"], columnKey: ["_field"], valueColumn: "_value")`, startTime, endTime))
	if err != nil {
		return nil, err
	}
	return events1, nil
}

// ReadingTotalCount returns the total count of Event from the database
func (c *HybridClient) ReadingTotalCount() (uint32, errors.EdgeX) {
	//Influx count
	count, err := TotalCountInflux(c,
		`
		|>range(start: 0, stop: now())
		|>group()
		|>filter(fn: (r) => r._field == "readingid")
		|>count()`)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// AllReadings query events by offset, limit, and labels
func (c *HybridClient) AllReadings(offset int, limit int) ([]model.Reading, errors.EdgeX) {
	readings, err := allInfluxReadings(c,
		fmt.Sprintf(`
		|>range(start: 0, stop: now())
		|>filter(fn: (r) => r["_field"] == "eventtype" or r["_field"] == "readingid" or r["_field"] == "units" or r["_field"] == "readingorigin" or r["_field"] == "valuetype" or r["_field"] == "value")
		|>group()
		|>sort(columns: ["_time"])
		|>pivot(rowKey: ["_time","counter","devicename","_measurement","resourcename"], columnKey: ["_field"], valueColumn: "_value")
		|>limit(n: %v, offset: %v)`, limit, offset))
	if err != nil {
		return nil, err
	}
	return readings, nil
}

// ReadingsByTimeRange query readings by time range, offset, and limit
func (c *HybridClient) ReadingsByTimeRange(start int, end int, offset int, limit int) (readings []model.Reading, edgeXerr errors.EdgeX) {
	readings, err := allInfluxReadings(c,
		fmt.Sprintf(`
		|>range(start:time(v:%v), stop: time(v:%v))
		|>filter(fn: (r) => r["_field"] == "eventtype" or r["_field"] == "readingid" or r["_field"] == "units" or r["_field"] == "readingorigin" or r["_field"] == "valuetype" or r["_field"] == "value")
		|>group()
		|>sort(columns: ["_time"])
		|>pivot(rowKey: ["_time","counter","devicename","_measurement","resourcename"], columnKey: ["_field"], valueColumn: "_value")
		|> limit(n: %v, offset: %v)`, start, end, limit, offset))
	if err != nil {
		return nil, err
	}
	return readings, nil
}

// ReadingsByResourceName query readings by offset, limit and resource name
func (c *HybridClient) ReadingsByResourceName(offset int, limit int, resourceName string) (readings []model.Reading, edgeXerr errors.EdgeX) {
	readings, err := allInfluxReadings(c,
		fmt.Sprintf(`
		|>range(start: 0, stop: now())
		|>filter(fn: (r) => r["resourcename"] == "%v")
		|>filter(fn: (r) => r["_field"] == "eventtype" or r["_field"] == "readingid" or r["_field"] == "units" or r["_field"] == "readingorigin" or r["_field"] == "valuetype" or r["_field"] == "value")
		|>group()
		|>sort(columns: ["_time"])
		|>pivot(rowKey: ["_time","counter","devicename","_measurement","resourcename"], columnKey: ["_field"], valueColumn: "_value")
		|>limit(n: %v, offset: %v)`, resourceName, limit, offset))
	if err != nil {
		return nil, err
	}
	return readings, nil
}

// ReadingsByDeviceName query readings by offset, limit and device name
func (c *HybridClient) ReadingsByDeviceName(offset int, limit int, name string) (readings []model.Reading, edgeXerr errors.EdgeX) {
	readings, err := allInfluxReadings(c,
		fmt.Sprintf(`
		|>range(start: 0, stop: now())
		|>filter(fn: (r) => r["devicename"] == "%v")
		|>filter(fn: (r) => r["_field"] == "eventtype" or r["_field"] == "readingid" or r["_field"] == "units" or r["_field"] == "readingorigin" or r["_field"] == "valuetype" or r["_field"] == "value")
		|>group()
		|>sort(columns: ["_time"])
		|>pivot(rowKey: ["_time","counter","devicename","_measurement","resourcename"], columnKey: ["_field"], valueColumn: "_value")
		|>limit(n: %v, offset: %v)`, name, limit, offset))
	if err != nil {
		return nil, err
	}
	return readings, nil
}

// ReadingCountByDeviceName returns the count of Readings associated a specific Device from the database
func (c *HybridClient) ReadingCountByDeviceName(deviceName string) (uint32, errors.EdgeX) {
	//Influx count
	count, err := TotalCountInflux(c,
		fmt.Sprintf(`
		|>range(start: 0, stop: now())
		|>filter(fn: (r) => r["devicename"] == "%v")
		|>group()
		|>filter(fn: (r) => r._field == "readingid")
		|>count()`, deviceName))
	if err != nil {
		return 0, err
	}

	return count, nil
}

// ReadingCountByResourceName returns the count of Readings associated a specific Resource from the database
func (c *HybridClient) ReadingCountByResourceName(resourceName string) (uint32, errors.EdgeX) {
	//Influx count
	count, err := TotalCountInflux(c,
		fmt.Sprintf(`
		|>range(start: 0, stop: now())
		|>filter(fn: (r) => r["resourcename"] == "%v")
		|>group()
		|>filter(fn: (r) => r._field == "readingid")
		|>count()`, resourceName))
	if err != nil {
		return 0, err
	}
	return count, nil
}

// ReadingCountByResourceNameAndTimeRange returns the count of Readings associated a specific Resource from the database within specified time range
func (c *HybridClient) ReadingCountByResourceNameAndTimeRange(resourceName string, startTime int, endTime int) (uint32, errors.EdgeX) {
	//Influx count
	count, err := TotalCountInflux(c,
		fmt.Sprintf(`
			|>range(start:time(v:%v), stop: time(v:%v))
			|>filter(fn: (r) => r["resourcename"] == "%v")
			|>group()
			|>filter(fn: (r) => r._field == "readingid")
			|>count()`, startTime, endTime, resourceName))
	if err != nil {
		return 0, err
	}
	return count, nil
}

// ReadingCountByDeviceNameAndResourceName returns the count of Readings associated with specified Resource and Device from the database
func (c *HybridClient) ReadingCountByDeviceNameAndResourceName(deviceName string, resourceName string) (uint32, errors.EdgeX) {
	//Influx count
	count, err := TotalCountInflux(c,
		fmt.Sprintf(`
			|>range(start: 0, stop: now())
			|>filter(fn: (r) => r["devicename"] == "%v" and r["resourcename"] == "%v" )
			|>group()
			|>filter(fn: (r) => r._field == "readingid")
			|>count()`, deviceName, resourceName))
	if err != nil {
		return 0, err
	}
	return count, nil
}

// ReadingCountByDeviceNameAndResourceNameAndTimeRange returns the count of Readings associated with specified Resource and Device from the database within specified time range
func (c *HybridClient) ReadingCountByDeviceNameAndResourceNameAndTimeRange(deviceName string, resourceName string, start int, end int) (uint32, errors.EdgeX) {
	//Influx count
	count, err := TotalCountInflux(c,
		fmt.Sprintf(`
			|>range(start:time(v:%v), stop: time(v:%v))
			|>filter(fn: (r) => r["devicename"] == "%v" and r["resourcename"] == "%v" )
			|>group()
			|>filter(fn: (r) => r._field == "readingid")
			|>count()`, start, end, deviceName, resourceName))
	if err != nil {
		return 0, err
	}
	return count, nil
}

// ReadingCountByTimeRange returns the count of Readings from the database within specified time range
func (c *HybridClient) ReadingCountByTimeRange(start int, end int) (uint32, errors.EdgeX) {
	//Influx count
	count, err := TotalCountInflux(c,
		fmt.Sprintf(`
			|>range(start:time(v:%v), stop: time(v:%v))			
			|>group()
			|>filter(fn: (r) => r._field == "readingid")
			|>count()`, start, end))
	if err != nil {
		return 0, err
	}
	return count, nil
}

// ReadingsByResourceNameAndTimeRange query readings by resourceName and specified time range. Readings are sorted in descending order of origin time.
func (c *HybridClient) ReadingsByResourceNameAndTimeRange(resourceName string, start int, end int, offset int, limit int) (readings []model.Reading, err errors.EdgeX) {
	readings, err = allInfluxReadings(c,
		fmt.Sprintf(`
		|>range(start: 0, stop: now())
		|>filter(fn: (r) => r["resourcename"] == "%v")
		|>filter(fn: (r) => r["_field"] == "eventtype" or r["_field"] == "readingid" or r["_field"] == "units" or r["_field"] == "readingorigin" or r["_field"] == "valuetype" or r["_field"] == "value")
		|>group()
		|>sort(columns: ["_time"])
		|>pivot(rowKey: ["_time","counter","devicename","_measurement","resourcename"], columnKey: ["_field"], valueColumn: "_value")
		|>limit(n: %v, offset: %v)`, resourceName, limit, offset))
	if err != nil {
		return nil, err
	}
	return readings, nil
}

func (c *HybridClient) ReadingsByDeviceNameAndResourceName(deviceName string, resourceName string, offset int, limit int) (readings []model.Reading, err errors.EdgeX) {
	readings, err = allInfluxReadings(c,
		fmt.Sprintf(`
		|>range(start: 0, stop: now())
		|>filter(fn: (r) => r["devicename"] == "%v" and r["resourcename"] == "%v")
		|>filter(fn: (r) => r["_field"] == "eventtype" or r["_field"] == "readingid" or r["_field"] == "units" or r["_field"] == "readingorigin" or r["_field"] == "valuetype" or r["_field"] == "value")
		|>group()
		|>sort(columns: ["_time"])
		|>pivot(rowKey: ["_time","counter","devicename","_measurement","resourcename"], columnKey: ["_field"], valueColumn: "_value")
		|>limit(n: %v, offset: %v)`, deviceName, resourceName, limit, offset))
	if err != nil {
		return nil, err
	}
	return readings, nil
}

func (c *HybridClient) ReadingsByDeviceNameAndResourceNameAndTimeRange(deviceName string, resourceName string, start int, end int, offset int, limit int) (readings []model.Reading, err errors.EdgeX) {
	readings, err = allInfluxReadings(c,
		fmt.Sprintf(`
		|>range(start:time(v:%v), stop: time(v:%v))	
		|>filter(fn: (r) => r["devicename"] == "%v" and r["resourcename"] == "%v")
		|>filter(fn: (r) => r["_field"] == "eventtype" or r["_field"] == "readingid" or r["_field"] == "units" or r["_field"] == "readingorigin" or r["_field"] == "valuetype" or r["_field"] == "value")
		|>group()
		|>sort(columns: ["_time"])
		|>pivot(rowKey: ["_time","counter","devicename","_measurement","resourcename"], columnKey: ["_field"], valueColumn: "_value")
		|>limit(n: %v, offset: %v)`, start, end, deviceName, resourceName, limit, offset))
	if err != nil {
		return nil, err
	}
	return readings, nil
}

func (c *HybridClient) ReadingsByDeviceNameAndResourceNamesAndTimeRange(deviceName string, resourceNames []string, start, end, offset, limit int) (readings []model.Reading, totalCount uint32, err errors.EdgeX) {
	var resourcePart string
	for _, resource := range resourceNames {
		resourcePart += fmt.Sprintf(` r["resourcename"] == "%v" or`, resource)
	}
	//remove the last or
	i := strings.LastIndex(resourcePart, "or")
	resourcePart = resourcePart[:i] + strings.Replace(resourcePart[i:], "or", "", 1)

	countchannel := make(chan uint32)
	readingchannel := make(chan []model.Reading)

	// the count execution is required as the array count of readings will consider the offset
	go func() {
		//Influx count
		count, err1 := TotalCountInflux(c,
			fmt.Sprintf(`
			|>range(start:time(v:%v), stop: time(v:%v))	
			|>filter(fn: (r) => r["devicename"] == "%v" and (%v))	
			|>group()
			|>filter(fn: (r) => r._field == "readingid")
			|>count()`, start, end, deviceName, resourcePart))
		if err1 != nil {
			err = err1
		}
		countchannel <- count
	}()

	go func() {
		readings, err1 := allInfluxReadings(c,
			fmt.Sprintf(`
			|>range(start:time(v:%v), stop: time(v:%v))
			|>filter(fn: (r) => r["devicename"] == "%v" and (%v))
			|>filter(fn: (r) => r["_field"] == "eventtype" or r["_field"] == "readingid" or r["_field"] == "units" or r["_field"] == "readingorigin" or r["_field"] == "valuetype" or r["_field"] == "value")
			|>group()
			|>sort(columns: ["_time"])
			|>pivot(rowKey: ["_time","counter","devicename","_measurement","resourcename"], columnKey: ["_field"], valueColumn: "_value")
			|>limit(n: %v, offset: %v)`, start, end, deviceName, resourcePart, limit, offset))
		if err1 != nil {
			err = err1
		}
		readingchannel <- readings
	}()
	//We wait for both the channels
	totalCount = <-countchannel
	readings = <-readingchannel

	if err != nil {
		return nil, 0, err
	}
	return readings, totalCount, nil
}

func (c *HybridClient) ReadingsByDeviceNameAndTimeRange(deviceName string, start int, end int, offset int, limit int) (readings []model.Reading, err errors.EdgeX) {
	readings, err = allInfluxReadings(c,
		fmt.Sprintf(`
		|>range(start:time(v:%v), stop: time(v:%v))	
		|>filter(fn: (r) => r["devicename"] == "%v")
		|>filter(fn: (r) => r["_field"] == "eventtype" or r["_field"] == "readingid" or r["_field"] == "units" or r["_field"] == "readingorigin" or r["_field"] == "valuetype" or r["_field"] == "value")
		|>group()
		|>sort(columns: ["_time"])
		|>pivot(rowKey: ["_time","counter","devicename","_measurement","resourcename"], columnKey: ["_field"], valueColumn: "_value")
		|>limit(n: %v, offset: %v)`, start, end, deviceName, limit, offset))
	if err != nil {
		return nil, err
	}
	return readings, nil
}

func (c *HybridClient) ReadingCountByDeviceNameAndTimeRange(deviceName string, start int, end int) (uint32, errors.EdgeX) {
	//Influx count
	count, err := TotalCountInflux(c,
		fmt.Sprintf(`
			|>range(start:time(v:%v), stop: time(v:%v))
			|>filter(fn: (r) => r["devicename"] == "%v")
			|>group()
			|>filter(fn: (r) => r._field == "readingid")
			|>count()`, start, end, deviceName))
	if err != nil {
		return 0, err
	}
	return count, nil
}
