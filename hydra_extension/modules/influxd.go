package modules

import (
	"encoding/json"
	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client/v2"
	"time"
	"yym/hydra_extension/hydra"
)

func init() {
	hydra.RegisterType("InfluxClient", newInfluxClient)
}

type InfluxConfig struct {
	Addr      string
	Username  string
	Password  string
	Timeout   time.Duration
	DB        string `codec:"db"`
	Precision string
}

func newInfluxClient(config InfluxConfig) (*InfluxClient, error) {
	if config.Precision == "" {
		config.Precision = "ns"
	}
	dbClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.Addr,
		Username: config.Username,
		Password: config.Password,
		Timeout:  config.Timeout})
	if err != nil {
		return nil, err
	}

	return &InfluxClient{client: dbClient, config: &config}, nil
}

type InfluxClient struct {
	client client.Client
	config *InfluxConfig
	cache  chan *InfluxPoint
}

type InfluxPoint struct {
	Name   string
	Tags   map[string]string
	Fields map[string]interface{}
	Time   time.Time
}

func (self *InfluxClient) AddPoint(point *InfluxPoint) error {
	config := client.BatchPointsConfig{Precision: self.config.Precision, Database: self.config.DB}
	bp, err := client.NewBatchPoints(config)
	if err != nil {
		return err
	}
	p, err := client.NewPoint(point.Name, point.Tags, point.Fields, point.Time)
	if err != nil {
		return err
	}
	bp.AddPoint(p)
	return self.client.Write(bp)
}

func (self *InfluxClient) AddPoints(points []*InfluxPoint) error {
	config := client.BatchPointsConfig{Precision: self.config.Precision, Database: self.config.DB}
	bp, err := client.NewBatchPoints(config)
	if err != nil {
		return err
	}
	for _, point := range points {
		p, err := client.NewPoint(point.Name, point.Tags, point.Fields, point.Time)
		if err != nil {
			return err
		}
		bp.AddPoint(p)
	}
	return self.client.Write(bp)
}

func (self *InfluxClient) Exec(sql string) error {
	resp, err := self.client.Query(client.NewQuery(sql, self.config.DB, "ns"))
	if err != nil {
		return err
	}
	if resp.Error() != nil {
		return resp.Error()
	}
	return nil
}

func (self *InfluxClient) Query(sql string, precision ...string) (interface{}, error) {
	p := ""
	if len(precision) > 0 {
		p = precision[0]
	}
	resp, err := self.client.Query(client.NewQuery(sql, self.config.DB, p))
	if err != nil {
		return nil, err
	}
	if resp.Error() != nil {
		return nil, resp.Error()
	}
	res := resp.Results
	ret := []map[string]interface{}{}
	if len(res) == 0 || len(res[0].Series) == 0 {
		return ret, nil
	}
	for _, row := range res[0].Series[0].Values {
		doc := map[string]interface{}{}
		for i, name := range res[0].Series[0].Columns {
			if i == 0 {
				if p == "no" {
					continue
				} else if p == "" {
					v, ok := row[i].(string)
					if ok {
						t, err := time.Parse(time.RFC3339, v)
						if t.Unix() != 0 && err == nil {
							doc[name] = t.Format("2006-01-02 15:04:05")
						}
					}
				} else if p == "s" {
					num, ok := row[i].(json.Number)
					if ok {
						doc[name] = num
					}
				} else if p == "ns" || p == "us" {
					num, ok := row[i].(json.Number)
					if ok {
						doc[name] = num.String()
					}
				} else {
					num, ok := row[i].(json.Number)
					if ok {
						doc[name] = num.String()
					}
				}

			} else {
				num, ok := row[i].(json.Number)
				if ok {
					f, _ := num.Float64()
					doc[name] = f
				} else {
					doc[name] = row[i]
				}
			}

		}
		ret = append(ret, doc)
	}

	return ret, nil
}
