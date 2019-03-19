package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gopkg.in/linkedin/goavro.v2"
	"os"
	"strconv"
	"strings"
	"time"
)

var mode="encoder"

type rowJSON struct {
	ts   time.Time
	VID  string `json:"vid"`
	PID  string `json:"pid"`
	MID  string `json:"mid"`
	ZID  string `json:"zid"`
	IP   string `json:"ip"`
	UID  string `json:"gid"`
	IDFA string `json:"idfa"`
	UA   string `json:"ua"`
	REF  string `json:"ref"`
	Lang string `json:"lang"`
	IID  string `json:"iid"`
	ATS  string `json:"ats"`
}

func newRowJSON(timeStr, str string) rowJSON {
	ret := rowJSON{}
	err := json.Unmarshal([]byte(str), &ret)
	if err != nil {
		//fmt.Println(err)
		return ret
	}
	ut, _ := strconv.ParseInt(timeStr, 10, 64)
	ret.ts = time.Unix(ut, 64)
	return ret
}

func (r rowJSON) valid() bool{
	return r.VID != "" &&
		r.PID != "" &&
		r.MID != "" &&
		r.ZID != ""
}

func (r rowJSON) mapString() map[string]interface{} {
	ret := map[string]interface{}{
		"time":                  r.ts,
		"external_publisher_id": r.VID,
		"external_media_id":     r.PID,
		"external_site_id":      r.MID,
		"external_sub_site_id":  r.ZID,
		"campaign_id":           "-",
		"session_id":            map[string]interface{}{"string": r.IID},
	}
	if !r.valid() {
		ret["failed"]=map[string]interface{}{}
	}
	if r.IP != "" {
		ret["ip_long"] = map[string]interface{}{"string": r.IP}
	}
	if r.UID != "" {
		ret["uid"] = map[string]interface{}{"string": r.UID}
	}
	if r.IDFA != "" {
		ret["device_id"] = map[string]interface{}{"string": r.IDFA}
	}
	if r.UA != "" {
		ret["user_agent"] = map[string]interface{}{"string": r.UA}
	}
	if r.Lang != "" {
		ret["language"] = map[string]interface{}{"string": r.Lang}
	}
	if r.ATS != "" {
		ut, err := strconv.ParseInt(r.ATS, 10, 64)
		if err == nil {
			ret["impression_time"] = map[string]interface{}{"long.timestamp-millis": time.Unix(ut, 0)}
		}
	}
	return ret
}

func (r rowJSON) jsonByte() []byte {
	s, _ := json.Marshal(r.mapString())
	return s
}

func encoder() {
	conf := goavro.OCFConfig{
		W: os.Stdout,
		Schema: `
            {
    "type": "record",
    "namespace":"geniee.co.jp",
    "name":"ssp",
    "fields":[
    {"name":"time",                 "type": {"type": "long", "logicalType": "timestamp-millis"}},
    {"name":"external_publisher_id","type": "string"},
    {"name":"external_media_id",    "type": "string"},
    {"name":"external_site_id",     "type": "string"},
    {"name":"external_sub_site_id", "type": "string"},
    {"name":"campaign_id",          "type": "string"},
    {"name":"conversion_type",      "type": ["null", "string"], "default": null},
    {"name":"attribution_time",     "type": ["null", {"type": "long", "logicalType": "timestamp-millis"}], "default": null},
    {"name":"impression_time",      "type": ["null", {"type": "long", "logicalType": "timestamp-millis"}], "default": null},
    {"name":"ip_long",              "type": ["null", "string"], "default": null},
    {"name":"uid",                  "type": ["null", "string"], "default": null},
    {"name":"device_id",            "type": ["null", "string"], "default": null},
    {"name":"user_agent",           "type": ["null", "string"], "default": null},
    {"name":"referer",              "type": ["null", "string"], "default": null},
    {"name":"platform",             "type": ["null", "string"], "default": null},
    {"name":"device_type",          "type": ["null", "string"], "default": null},
    {"name":"os_version",           "type": ["null", "string"], "default": null},
    {"name":"language",             "type": ["null", "string"], "default": null},
    {"name":"wifi",                 "type": ["null", "string"], "default": null},
    {"name":"session_id",           "type": ["null", "string"], "default": null}
]
            }`,
	}
	writer, err := goavro.NewOCFWriter(conf)
	if err != nil {
		fmt.Println(err)
		return
	}
	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		tsv := strings.Split(stdin.Text(), "\t")
		if len(tsv) != 3 {
			continue
		}
		jsonMaps := []map[string]interface{}{newRowJSON(tsv[0], tsv[2]).mapString()}
		if _,exist :=jsonMaps[0]["failed"]; exist{
			//fmt.Println("failed")
			continue
		}
		err = writer.Append(jsonMaps)
		if err != nil {
			continue
		}
	}
}

func decoder() {
	reader, err := goavro.NewOCFReader(os.Stdin)
	if err != nil {
		fmt.Println(err)
	}
	for reader.Scan() {
		hoge, err := reader.Read()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(hoge)
	}
}
func main() {
	if mode=="encoder" {
		encoder()
	} else {
		decoder()
	}
}
