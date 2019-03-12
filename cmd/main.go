package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"encoding/json"
	"gopkg.in/linkedin/goavro.v2"
)

type rowJSON struct{
	ts string
	VID string `json:"vid"`
	PID string `json:"pid"`
	MID string `json:"mid"`
	ZID string `json:"zid"`
	IP string `json:"ip"`
	UID string `json:"gid"`
	IDFA string `json:"idfa"`
	UA string `  json:"ua"`
	REF string  `json:"ref"`
	Lang string `json:"lang"`
	IID string   `json:"iid"`
}

func newRowJSON(time,str string) rowJSON{
	ret := rowJSON{}
	json.Unmarshal([]byte(str),&ret)
	ret.ts=time
	return ret
}

type phybbitJSON struct{
	Time string `json:"time"`
	VID string `json:"external_publisher_id"`
	PID string `json:"external_media_id"`
	MID string `json:"external_site_id"`
	ZID string `json:"external_sub_site_id"`
	CID string `json:"campaign_id"`
	IP string `json:"ip_long"`
	UID string `json:"uid"`
	IDFA string `json:"device_id"`
	UA string `  json:"user_agent"`
	REF string  `json:"referer"`
	Lang string `json:"language"`
	IID string   `json:"session_id"`

}

func (r rowJSON)newPhybbitJSON() phybbitJSON{
	return phybbitJSON{
		Time: r.ts, 
		VID :r.VID,
		PID :r.PID,
		MID :r.MID,
		ZID :r.ZID,
		CID : "-",
		IP:r.IP,
		UID: r.UID,
		IDFA:r.IDFA,
		UA :r.UA,
		REF :r.REF,
		Lang:r.Lang,
		IID :r.IID,
	}
}

func (r phybbitJSON)bytes() ([]byte,bool,error){
	if r.ZID == "1306008" {
		return nil,false,nil
	}
	b, err := json.Marshal(r)
	if err != nil {
		return nil, true,err
	}
	return b, true,nil
}

func main(){
	codec, err := goavro.NewCodec(
		`
            {
    "type": "record",

    "name":"ssp",
    "fields":[
    {"name":"time",                 "type": "string"},
    {"name":"external_publisher_id","type": "string"},
    {"name":"external_media_id",    "type": "string"},
    {"name":"external_site_id",     "type": "string"},
    {"name":"external_sub_site_id", "type": "string"},
    {"name":"campaign_id",          "type": "string"},
    {"name":"conversion_type",      "type": ["null", "string"], "default": null},
    {"name":"attribution_time",     "type": ["null", "string"], "default": null},
    {"name":"impression_time",      "type": ["null", "string"], "default": null},
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
            }`)
	if err != nil {
		fmt.Println(err)
	}
	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan(){
		line := stdin.Text()
		tsv:=strings.Split(line,"\t")
		jsonByte,valid,err:=newRowJSON(tsv[0], tsv[2]).newPhybbitJSON().bytes()
		if !valid {
			continue
		}
		if err != nil{
			fmt.Println(err)
			continue
		}

		native,_,err := codec.NativeFromTextual(jsonByte)
		if err != nil {
			fmt.Println(err)
			continue
		}
		binary, err := codec.BinaryFromNative(nil, native)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(binary)
	}
}
