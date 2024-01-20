//
// SPDX-License-Identifier: Apache-2.0
//
// Client wrapping most Influx APIs

package influx

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"crypto/sha512"
	b64 "encoding/base64"

	"github.com/edgexfoundry/edgex-go/internal/pkg/db"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/errors"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type Client struct {
	influxClient influxdb2.Client
	influxOrg    string
	influxBucket string
	batchSize    uint64
}

func NewClient(config db.Configuration, lc logger.LoggingClient) (*Client, errors.EdgeX) {

	influxUrl := os.Getenv("INFLUXDB_URL")
	if influxUrl == "" {
		return nil, errors.NewCommonEdgeX(errors.KindDatabaseError, "INFLUXDB_URL must be set", nil)
	}
	influxOrg := os.Getenv("INFLUXDB_ORG")
	if influxOrg == "" {
		return nil, errors.NewCommonEdgeX(errors.KindDatabaseError, "INFLUXDB_ORG must be set", nil)
	}
	influxBucket := os.Getenv("INFLUXDB_BUCKET")
	if influxBucket == "" {
		return nil, errors.NewCommonEdgeX(errors.KindDatabaseError, "INFLUXDB_BUCKET must be set", nil)
	}

	var user, pass string
	if os.Getenv("EDGEX_SECURITY_SECRET_STORE") != "false" {
		user = "admin"
		pass = config.Password

	} else {
		user = "admin"
		pass = "admin1234"
	}
	//Create influx token
	h := sha512.New()
	h.Write([]byte(user + pass))
	bs := h.Sum(nil)
	influxToken := b64.URLEncoding.EncodeToString(bs)
	var influxClient influxdb2.Client
	influxClient = influxdb2.NewClient(influxUrl, influxToken)
	//See if organization exists
	_, err := influxClient.OrganizationsAPI().FindOrganizationByName(context.Background(), influxOrg)
	if err != nil {
		//Try to do the setup of Influx
		influxClient = influxdb2.NewClient(influxUrl, "")
		_, err = influxClient.SetupWithToken(context.Background(), user, pass, influxOrg, "Edgex", 0, influxToken)
		if err != nil {
			return nil, errors.NewCommonEdgeX(errors.KindDatabaseError, "INFLUXDB setup failed", nil)
		}
	}
	var batch uint64
	influxBatch := os.Getenv("INFLUXDB_BATCH")
	if influxBatch == "" {
		influxClient = influxdb2.NewClient(influxUrl, influxToken)
	} else {
		batch, _ = strconv.ParseUint(influxBatch, 10, 32)
		influxClient = influxdb2.NewClientWithOptions(influxUrl, influxToken,
			influxdb2.DefaultOptions().SetBatchSize(uint(batch)))
	}
	// validate client connection health
	_, err = influxClient.Health(context.Background())
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindDatabaseError, "INFLUXDB cannot be reached", nil)
	}

	influx := &Client{influxClient: influxClient,
		influxOrg: influxOrg, influxBucket: influxBucket, batchSize: batch}

	return influx, nil
}

// CreateBucket creates a bucket in InfluxDB if it does not already exist.
// An existing bucket will be reused
func (c *Client) CreateBucket() errors.EdgeX {

	ctx := context.Background()
	bucketsAPI := c.influxClient.BucketsAPI()

	//Find the bucket with ord id
	domainOrg, _ := c.influxClient.OrganizationsAPI().FindOrganizationByName(ctx, c.influxOrg)
	bucketList, err := bucketsAPI.FindBucketsByOrgID(ctx, *domainOrg.Id)
	if err != nil {
		return errors.NewCommonEdgeX(errors.KindInvalidId, "Bucket list cannot be retrieved", err)
	}
	for _, bucket := range *bucketList {
		if bucket.Name == c.influxBucket {
			return nil
		}
	}
	// create new empty bucket
	_, err = c.influxClient.BucketsAPI().CreateBucketWithNameWithID(ctx, *domainOrg.Id, c.influxBucket)

	if err != nil {
		return errors.NewCommonEdgeX(errors.KindInvalidId, "Bucket creation failed", err)
	}
	return nil
}

func (c *Client) WritePoint(measurement string, tags map[string]string, values map[string]interface{}, time time.Time) {
	writeAPI := c.influxClient.WriteAPI(c.influxOrg, c.influxBucket)
	p := influxdb2.NewPoint(measurement, tags, values, time)
	writeAPI.WritePoint(p)
}

func (c *Client) QueryData(query string) (*api.QueryTableResult, error) {

	//Add the bucket information in query
	query = fmt.Sprintf(`from(bucket: "%v")%v`, c.influxBucket, query)
	// Get query client
	queryAPI := c.influxClient.QueryAPI(c.influxOrg)
	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	return result, nil

}

func (c *Client) DeleteData(deletecondition string, starttime time.Time, endtime time.Time) error {
	// Get delete client
	deleteAPI := c.influxClient.DeleteAPI()
	return deleteAPI.DeleteWithName(context.Background(), c.influxOrg, c.influxBucket, starttime, endtime, deletecondition)
}

func (c *Client) CloseSession() {
	c.influxClient.Close()
}
