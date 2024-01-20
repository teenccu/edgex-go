//
// SPDX-License-Identifier: Apache-2.0
//
// Influx code handling both edgex events,readings and influx APIs for traversing records

package hybrid

import (
	"fmt"
	"strconv"
	"time"

	"encoding/hex"
	"encoding/json"

	"github.com/edgexfoundry/go-mod-core-contracts/v3/errors"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/models"
	"github.com/google/uuid"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

func TotalCountInflux(conn *HybridClient, fluxQuery string) (uint32, errors.EdgeX) {
	var count uint32
	result, err := conn.influxClient.QueryData(fluxQuery)
	if err == nil {
		// Iterate over query response
		for result.Next() {
			c, _ := strconv.ParseUint(fmt.Sprintf("%v", result.Record().Value()), 10, 32)
			count = uint32(c)
			break
		}
	} else {
		return 0, errors.NewCommonEdgeX(errors.KindDatabaseError, "INFLUXDB count cannnot be retreived", err)
	}

	return count, nil
}

func createReading(result *api.QueryTableResult) (newrdng models.Reading) {
	//Create the simple reading object
	if result.Record().ValueByKey("eventtype") == "sr" {
		newrdng := models.SimpleReading{}
		newrdng.Id = fmt.Sprintf("%v", result.Record().ValueByKey("readingid"))
		newrdng.ResourceName = fmt.Sprintf("%v", result.Record().ValueByKey("resourcename"))
		newrdng.DeviceName = fmt.Sprintf("%v", result.Record().ValueByKey("devicename"))
		newrdng.ProfileName = fmt.Sprintf("%v", result.Record().Measurement())
		newrdng.ValueType = fmt.Sprintf("%v", result.Record().ValueByKey("valuetype"))
		newrdng.Units = fmt.Sprintf("%v", result.Record().ValueByKey("units"))
		if newrdng.Units == "<nil>" {
			newrdng.Units = ""
		}

		newrdng.Value = fmt.Sprintf("%v", result.Record().ValueByKey("value"))
		newrdng.Origin, _ = strconv.ParseInt(fmt.Sprintf("%v", result.Record().ValueByKey("readingorigin")), 10, 64)
		newrdng.Tags = map[string]any{} //TODO
		return newrdng
	}

	//Create the binary reading object
	if result.Record().ValueByKey("eventtype") == "br" {
		newrdng := models.BinaryReading{}
		newrdng.Id = fmt.Sprintf("%v", result.Record().ValueByKey("readingid"))
		newrdng.ResourceName = fmt.Sprintf("%v", result.Record().ValueByKey("resourcename"))
		newrdng.DeviceName = fmt.Sprintf("%v", result.Record().ValueByKey("devicename"))
		newrdng.ProfileName = fmt.Sprintf("%v", result.Record().Measurement())
		newrdng.ValueType = fmt.Sprintf("%v", result.Record().ValueByKey("valuetype"))
		newrdng.Units = fmt.Sprintf("%v", result.Record().ValueByKey("units"))
		if newrdng.Units == "<nil>" {
			newrdng.Units = ""
		}
		newrdng.MediaType = fmt.Sprintf("%v", result.Record().ValueByKey("mediatype"))
		newrdng.BinaryValue, _ = hex.DecodeString(fmt.Sprintf("%v", result.Record().ValueByKey("value")))
		newrdng.Origin, _ = strconv.ParseInt(fmt.Sprintf("%v", result.Record().ValueByKey("readingorigin")), 10, 64)
		newrdng.Tags = map[string]any{} //TODO
		return newrdng
	}

	//Create the object reading object
	if result.Record().ValueByKey("eventtype") == "or" {
		newrdng := models.ObjectReading{}
		newrdng.Id = fmt.Sprintf("%v", result.Record().ValueByKey("readingid"))
		newrdng.ResourceName = fmt.Sprintf("%v", result.Record().ValueByKey("resourcename"))
		newrdng.DeviceName = fmt.Sprintf("%v", result.Record().ValueByKey("devicename"))
		newrdng.ProfileName = fmt.Sprintf("%v", result.Record().Measurement())
		newrdng.ValueType = fmt.Sprintf("%v", result.Record().ValueByKey("valuetype"))
		newrdng.Units = fmt.Sprintf("%v", result.Record().ValueByKey("units"))
		if newrdng.Units == "<nil>" {
			newrdng.Units = ""
		}
		hexarry, _ := hex.DecodeString(fmt.Sprintf("%v", result.Record().ValueByKey("value")))
		var obj interface{}
		_ = json.Unmarshal(hexarry, &obj)
		newrdng.ObjectValue = obj
		newrdng.Origin, _ = strconv.ParseInt(fmt.Sprintf("%v", result.Record().ValueByKey("readingorigin")), 10, 64)
		newrdng.Tags = map[string]any{} //TODO
		return newrdng
	}
	return nil
}

func allInfluxReadings(conn *HybridClient, fluxQuery string) ([]models.Reading, errors.EdgeX) {
	resultsArr := []models.Reading{}
	var reading models.Reading
	result, err := conn.influxClient.QueryData(fluxQuery)
	if err == nil {
		// Iterate over query response
		for result.Next() {

			// check for an error
			if result.Err() != nil {
				fmt.Printf("query parsing error: %s\n", result.Err().Error())
			}
			reading = createReading(result)
			resultsArr = append(resultsArr, reading)
		}

	} else {
		return nil, errors.NewCommonEdgeX(errors.KindDatabaseError, "INFLUXDB reading error", err)
	}

	return resultsArr, nil
}

func AllEvents(conn *HybridClient, offset int, limit int, fluxQuery string) ([]models.Event, errors.EdgeX) {

	if offset <= 0 {
		offset = 1
	}

	if limit == 0 {
		limit = -2 // limit is ineffective
	}
	resultsArr := []models.Event{}
	currentevent := models.Event{}
	var curid string
	var reading models.Reading
	result, err := conn.influxClient.QueryData(fluxQuery)
	if err == nil {
		// Iterate over query response
		for result.Next() {

			// check for an error
			if result.Err() != nil {
				fmt.Printf("query parsing error: %s\n", result.Err().Error())
			}

			// Detect if new event
			if curid == result.Record().ValueByKey("eventid").(string) {
				if offset > 1 {
					continue // Do not consider reading in offset calcul
				}
				//Create the reading object
				reading = createReading(result)
				currentevent.Readings = append(currentevent.Readings, reading)
			} else {
				// Add the last event
				if currentevent.Id != "" { // Ignores the first one
					resultsArr = append(resultsArr, currentevent)
				}
				// Check offset and limit when there is a new event
				if offset > 1 {
					offset--
					continue
				} else {
					limit--
					if limit == -1 {
						return resultsArr, nil
					}
				}
				curid = result.Record().ValueByKey("eventid").(string)
				currentevent = models.Event{
					Id:          fmt.Sprintf("%v", result.Record().ValueByKey("eventid")),
					DeviceName:  fmt.Sprintf("%v", result.Record().ValueByKey("devicename")),
					ProfileName: fmt.Sprintf("%v", result.Record().Measurement()),
					SourceName:  fmt.Sprintf("%v", result.Record().ValueByKey("sourcename")),
					Origin:      result.Record().Time().UnixNano(),
				}

				currentevent.Readings = []models.Reading{}
				//Create the reading object
				reading = createReading(result)
				currentevent.Readings = append(currentevent.Readings, reading)
			}
		}

	} else {
		return nil, errors.NewCommonEdgeX(errors.KindDatabaseError, "INFLUXDB allevents error", err)
	}

	// Add any event that is created
	if currentevent.Id != "" {
		resultsArr = append(resultsArr, currentevent)
	}
	return resultsArr, nil
}

func AddEvent(conn *HybridClient, e models.Event) (addedEvent models.Event, edgeXerr errors.EdgeX) {
	var baseReading *models.BaseReading
	var err error
	unixtime := time.Unix(0, int64(e.Origin))
	var counter int
	for _, r := range e.Readings {
		counter++
		switch newReading := r.(type) {

		case models.BinaryReading: //https://docs.edgexfoundry.org/3.0/examples/Ch-ExamplesSendingAndConsumingBinary/
			baseReading = &newReading.BaseReading
			if err = checkReadingValue(baseReading); err != nil {
				return e, errors.NewCommonEdgeXWrapper(err)
			}
			conn.influxClient.WritePoint(e.ProfileName,
				map[string]string{"devicename": e.DeviceName, "resourcename": r.GetBaseReading().ResourceName, "counter": strconv.Itoa(counter)}, // this is decided as most of queries in the interface are by device name and then resourcename. Also tried to reduce high cardinality
				map[string]interface{}{"eventid": e.Id, "readingid": r.GetBaseReading().Id, "readingorigin": fmt.Sprintf("%v", r.GetBaseReading().Origin), "eventtype": "br", "units": newReading.Units, "valuetype": newReading.ValueType, "mediatype": newReading.MediaType, "value": hex.EncodeToString(newReading.BinaryValue)},
				unixtime)

		case models.SimpleReading:
			baseReading = &newReading.BaseReading
			if err = checkReadingValue(baseReading); err != nil {
				return e, errors.NewCommonEdgeXWrapper(err)
			}

			conn.influxClient.WritePoint(e.ProfileName,
				map[string]string{"devicename": e.DeviceName, "resourcename": r.GetBaseReading().ResourceName, "counter": strconv.Itoa(counter)}, // this is decided as most of queries in the interface are by device name and then resourcename. Also tried to reduce high cardinality
				map[string]interface{}{"eventid": e.Id, "readingid": r.GetBaseReading().Id, "readingorigin": fmt.Sprintf("%v", r.GetBaseReading().Origin), "eventtype": "sr", "sourcename": e.SourceName, "units": newReading.Units, "valuetype": newReading.ValueType, "value": newReading.Value},
				unixtime)

		case models.ObjectReading:
			baseReading = &newReading.BaseReading
			if err = checkReadingValue(baseReading); err != nil {
				return e, errors.NewCommonEdgeXWrapper(err)
			}
			s, err := json.Marshal(newReading.ObjectValue)
			if err != nil {
				return e, errors.NewCommonEdgeXWrapper(err)
			}

			conn.influxClient.WritePoint(e.ProfileName,
				map[string]string{"devicename": e.DeviceName, "resourcename": r.GetBaseReading().ResourceName, "counter": strconv.Itoa(counter)}, // this is decided as most of queries in the interface are by device name and then resourcename. Also tried to reduce high cardinality
				map[string]interface{}{"eventid": e.Id, "readingid": r.GetBaseReading().Id, "readingorigin": fmt.Sprintf("%v", r.GetBaseReading().Origin), "eventtype": "or", "sourcename": e.SourceName, "units": newReading.Units, "valuetype": newReading.ValueType, "value": hex.EncodeToString(s)},
				unixtime)

		default:
			return e, errors.NewCommonEdgeX(errors.KindContractInvalid, "unsupported reading type", nil)
		}
	}
	return e, nil // TODO sort the reading as per resources as done in add events
}

func checkReadingValue(b *models.BaseReading) errors.EdgeX {
	// check if id is a valid uuid
	if b.Id == "" {
		b.Id = uuid.New().String()
	} else {
		_, err := uuid.Parse(b.Id)
		if err != nil {
			return errors.NewCommonEdgeX(errors.KindInvalidId, "uuid parsing failed", err)
		}
	}
	return nil
}
