package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"text/template"
	"time"
)

var (
	fleetQuery, hostsQuery *template.Template
)

func init() {
	fleetQuery = template.Must(template.New("fleetQuery").Parse(rawFleetQuery))
	hostsQuery = template.Must(template.New("hostsQuery").Parse(rawHostsQuery))

}

// ElasticSearchDevops produces ES-specific queries for the devops use case.
type ElasticSearchDevops struct {
	AllInterval TimeInterval
}

// NewElasticSearchDevops makes an ElasticSearchDevops object ready to generate Queries.
func NewElasticSearchDevops(start, end time.Time) *ElasticSearchDevops {
	if !start.Before(end) {
		panic("bad time order")
	}
	return &ElasticSearchDevops{
		AllInterval: NewTimeInterval(start, end),
	}
}

// Dispatch fulfills the QueryGenerator interface.
func (d *ElasticSearchDevops) Dispatch(i int, q *Query, scaleVar int) {
	DevopsDispatch(d, i, q, scaleVar)
}

// AvgCPUUsageDayByHour populates a Query for getting the average CPU usage by hour for a random day.
func (d *ElasticSearchDevops) AvgCPUUsageDayByHour(q *Query) {
	interval := d.AllInterval.RandWindow(24 * time.Hour)

	body := new(bytes.Buffer)
	mustExecuteTemplate(fleetQuery, body, FleetQueryParams{
		Start:  interval.StartString(),
		End:    interval.EndString(),
		Bucket: "1h",
		Field:  "usage_user",
	})

	humanLabel := []byte("Elastic avg cpu, all hosts, rand 1d by 1h")
	q.HumanLabel = humanLabel
	q.HumanDescription = []byte(fmt.Sprintf("%s: %s", humanLabel, interval.StartString()))
	q.Method = []byte("POST")

	q.Path = []byte("/cpu/_search")
	q.Body = body.Bytes()
}

// AvgCPUUsageWeekByHour populates a Query for getting the average CPU usage by hour for a random week.
func (d *ElasticSearchDevops) AvgCPUUsageWeekByHour(q *Query) {
	interval := d.AllInterval.RandWindow(7 * 24 * time.Hour)

	body := new(bytes.Buffer)
	mustExecuteTemplate(fleetQuery, body, FleetQueryParams{
		Start:  interval.StartString(),
		End:    interval.EndString(),
		Bucket: "1h",
		Field:  "usage_user",
	})

	humanLabel := []byte("Elastic avg cpu, all hosts, rand 7d by 1h")
	q.HumanLabel = humanLabel
	q.HumanDescription = []byte(fmt.Sprintf("%s: %s", humanLabel, interval.StartString()))
	q.Method = []byte("POST")

	q.Path = []byte("/cpu/_search")
	q.Body = body.Bytes()
}

// AvgCPUUsageMonthByDay populates a Query for getting the average CPU usage by day for a random month.
func (d *ElasticSearchDevops) AvgCPUUsageMonthByDay(q *Query) {
	interval := d.AllInterval.RandWindow(28 * 24 * time.Hour)

	body := new(bytes.Buffer)
	mustExecuteTemplate(fleetQuery, body, FleetQueryParams{
		Start:  interval.StartString(),
		End:    interval.EndString(),
		Bucket: "1d",
		Field:  "usage_user",
	})

	humanLabel := []byte("Elastic avg cpu, all hosts, rand 28d by 1d")
	q.HumanLabel = humanLabel
	q.HumanDescription = []byte(fmt.Sprintf("%s: %s", humanLabel, interval.StartString()))
	q.Method = []byte("POST")

	q.Path = []byte("/cpu/_search")
	q.Body = body.Bytes()
}

// AvgMemAvailableDayByHour populates a Query for getting the average memory available by hour for a random day.
func (d *ElasticSearchDevops) AvgMemAvailableDayByHour(q *Query) {
	interval := d.AllInterval.RandWindow(24 * time.Hour)

	body := new(bytes.Buffer)
	mustExecuteTemplate(fleetQuery, body, FleetQueryParams{
		Start:  interval.StartString(),
		End:    interval.EndString(),
		Bucket: "1h",
		Field:  "available",
	})

	humanLabel := []byte("Elastic avg mem, all hosts, rand 1d by 1h")
	q.HumanLabel = humanLabel
	q.HumanDescription = []byte(fmt.Sprintf("%s: %s", humanLabel, interval.StartString()))
	q.Method = []byte("POST")

	q.Path = []byte("/mem/_search")
	q.Body = body.Bytes()
}

// AvgMemAvailableWeekByHour populates a Query for getting the average memory available by hour for a random week.
func (d *ElasticSearchDevops) AvgMemAvailableWeekByHour(q *Query) {
	interval := d.AllInterval.RandWindow(7 * 24 * time.Hour)

	body := new(bytes.Buffer)
	mustExecuteTemplate(fleetQuery, body, FleetQueryParams{
		Start:  interval.StartString(),
		End:    interval.EndString(),
		Bucket: "1h",
		Field:  "available",
	})

	humanLabel := []byte("Elastic avg mem, all hosts, rand 7d by 1h")
	q.HumanLabel = humanLabel
	q.HumanDescription = []byte(fmt.Sprintf("%s: %s", humanLabel, interval.StartString()))
	q.Method = []byte("POST")

	q.Path = []byte("/mem/_search")
	q.Body = body.Bytes()
}

// AvgMemAvailableMonthByDay populates a Query for getting the average memory available by day for a random month.
func (d *ElasticSearchDevops) AvgMemAvailableMonthByDay(q *Query) {
	interval := d.AllInterval.RandWindow(28 * 24 * time.Hour)

	body := new(bytes.Buffer)
	mustExecuteTemplate(fleetQuery, body, FleetQueryParams{
		Start:  interval.StartString(),
		End:    interval.EndString(),
		Bucket: "1d",
		Field:  "available",
	})

	humanLabel := []byte("Elastic avg mem, all hosts, rand 28d by 1d")
	q.HumanLabel = humanLabel
	q.HumanDescription = []byte(fmt.Sprintf("%s: %s", humanLabel, interval.StartString()))
	q.Method = []byte("POST")

	q.Path = []byte("/mem/_search")
	q.Body = body.Bytes()
}

// MaxCPUUsageHourByMinuteOneHost populates a Query for getting the maximum CPU
// usage for one host over the course of an hour.
func (d *ElasticSearchDevops) MaxCPUUsageHourByMinuteOneHost(q *Query, scaleVar int) {
	d.maxCPUUsageHourByMinuteNHosts(q, scaleVar, 1)
}

// MaxCPUUsageHourByMinuteTwoHosts populates a Query for getting the maximum CPU
// usage for two hosts over the course of an hour.
func (d *ElasticSearchDevops) MaxCPUUsageHourByMinuteTwoHosts(q *Query, scaleVar int) {
	d.maxCPUUsageHourByMinuteNHosts(q, scaleVar, 2)
}

// MaxCPUUsageHourByMinuteFourHosts populates a Query for getting the maximum CPU
// usage for four hosts over the course of an hour.
func (d *ElasticSearchDevops) MaxCPUUsageHourByMinuteFourHosts(q *Query, scaleVar int) {
	d.maxCPUUsageHourByMinuteNHosts(q, scaleVar, 4)
}

// MaxCPUUsageHourByMinuteEightHosts populates a Query for getting the maximum CPU
// usage for four hosts over the course of an hour.
func (d *ElasticSearchDevops) MaxCPUUsageHourByMinuteEightHosts(q *Query, scaleVar int) {
	d.maxCPUUsageHourByMinuteNHosts(q, scaleVar, 8)
}

// MaxCPUUsageHourByMinuteSixteenHosts populates a Query for getting the maximum CPU
// usage for four hosts over the course of an hour.
func (d *ElasticSearchDevops) MaxCPUUsageHourByMinuteSixteenHosts(q *Query, scaleVar int) {
	d.maxCPUUsageHourByMinuteNHosts(q, scaleVar, 16)
}

func (d *ElasticSearchDevops) MaxCPUUsageHourByMinuteThirtyTwoHosts(q *Query, scaleVar int) {
	d.maxCPUUsageHourByMinuteNHosts(q, scaleVar, 32)
}

func (d *ElasticSearchDevops) maxCPUUsageHourByMinuteNHosts(q *Query, scaleVar, nhosts int) {
	interval := d.AllInterval.RandWindow(time.Hour)
	nn := rand.Perm(scaleVar)[:nhosts]

	hostnames := []string{}
	for _, n := range nn {
		hostnames = append(hostnames, fmt.Sprintf("host_%d", n))
	}

	hostnameClauses := []string{}
	for _, s := range hostnames {
		hostnameClauses = append(hostnameClauses, fmt.Sprintf("\"%s\"", s))
	}

	combinedHostnameClause := fmt.Sprintf("[ %s ]", strings.Join(hostnameClauses, ", "))

	body := new(bytes.Buffer)
	mustExecuteTemplate(hostsQuery, body, HostsQueryParams{
		JSONEncodedHostnames: combinedHostnameClause,
		Start:                interval.StartString(),
		End:                  interval.EndString(),
		Bucket:               "1m",
		Field:                "usage_user",
	})

	humanLabel := []byte(fmt.Sprintf("Elastic max cpu, rand %4d hosts, rand 1hr by 1m", nhosts))
	q.HumanLabel = humanLabel
	q.HumanDescription = []byte(fmt.Sprintf("%s: %s", humanLabel, interval.StartString()))
	q.Method = []byte("POST")

	q.Path = []byte("/cpu/_search")
	q.Body = body.Bytes()
}

func mustExecuteTemplate(t *template.Template, w io.Writer, params interface{}) {
	err := t.Execute(w, params)
	if err != nil {
		panic(fmt.Sprintf("logic error in executing template: %s", err))
	}
}

type FleetQueryParams struct {
	Bucket, Start, End, Field string
}

type HostsQueryParams struct {
	JSONEncodedHostnames      string
	Bucket, Start, End, Field string
}

const rawFleetQuery = `
{
  "size" : 0,
  "aggs": {
    "result": {
      "filter": {
        "range": {
          "timestamp": {
            "gte": "{{.Start}}",
            "lt": "{{.End}}"
          }
        }
      },
      "aggs": {
        "result2": {
          "date_histogram": {
            "field": "timestamp",
            "interval": "{{.Bucket}}",
            "format": "yyyy-MM-dd-HH"
          },
          "aggs": {
            "avg_of_field": {
              "avg": {
                 "field": "{{.Field}}"
              }
            }
          }
        }
      }
    }
  }
}
`

const rawHostsQuery = `
{
  "size":0,
  "aggs":{
    "result":{
      "filter":{
        "bool":{
          "filter":{
            "range":{
              "timestamp":{
                "gte":"{{.Start}}",
                "lt":"{{.End}}"
              }
            }
          },
          "should":[
            {
              "terms":{
                "hostname": {{.JSONEncodedHostnames }}
              }
            }
          ],
	  "minimum_should_match" : 1
        }
      },
      "aggs":{
        "result2":{
          "date_histogram":{
            "field":"timestamp",
            "interval":"{{.Bucket}}",
            "format":"yyyy-MM-dd-HH"
          },
          "aggs":{
            "max_of_field":{
              "max":{
                "field":"{{.Field}}"
              }
            }
          }
        }
      }
    }
  }
}
`
