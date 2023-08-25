package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/baowk/dilu-core/common/utils/http_util"
	"github.com/gin-gonic/gin"
)

func GetIP(c *gin.Context) string {
	ip := c.Request.Header.Get("X-Forwarded-For")
	if strings.Contains(ip, "127.0.0.1") || ip == "" {
		ip = c.Request.Header.Get("X-real-ip")
	}
	if ip == "" {
		ip = "127.0.0.1"
	}
	RemoteIP := c.RemoteIP()
	if RemoteIP != "127.0.0.1" {
		ip = RemoteIP
	}
	ClientIP := c.ClientIP()
	if ClientIP != "127.0.0.1" {
		ip = ClientIP
	}
	return ip
}

type IPLocation struct {
	Code int            `json:"code"` //返回码 200成功
	Msg  string         `json:"msg"`  //返回消息
	Data IPLocationData `json:"data"`
}

type IPLocationData struct {
	AreaCode       string   `json:"area_code"`       //: "320311",
	Province       string   `json:"province"`        //省: "江苏",
	City           string   `json:"city"`            //: "徐州",
	District       string   `json:"district"`        //: "丰县",
	CityCode       string   `json:"city_code"`       //: "0516",
	Continent      string   `json:"continent"`       //: "亚洲",
	Country        string   `json:"country"`         //: "中国",
	CountryCode    string   `json:"country_code"`    //: "CN",
	CountryEnglish string   `json:"country_english"` //: "",
	Elevation      string   `json:"elevation"`       //: "40",
	Ip             string   `json:"ip"`              //: "114.234.76.140",
	Isp            string   `json:"isp"`             //: "电信",
	Latitude       string   `json:"latitude"`        //: "34.227883",
	LocalTime      string   `json:"local_time"`      //: "2023-08-02 14:36",
	Longitude      string   `json:"longitude"`       //: "117.213995",
	MultiStreet    []Street `json:"multi_street"`
	Street         string   `json:"street"`          //: "解放路168号",
	Version        string   `json:"version"`         //: "V4",
	WeatherStation string   `json:"weather_station"` //: "CHXX0437",
	ZipCode        string   `json:"zip_code"`        //: "221006"
}

type Street struct {
	Lng          string `json:"lng"`           //经度: "116.60833",
	Lat          string `json:"lat"`           //纬度: "34.701533",
	Province     string `json:"province"`      //省: "江苏",
	City         string `json:"city"`          //: "徐州",
	District     string `json:"district"`      //: "丰县",
	Street       string `json:"street"`        //: "解放路168号",
	StreetNumber string `json:"street_number"` //: "解放路168号"
}

func GetLocationByIp(secretKey, ip string, location *IPLocationData) error {
	url := "https://api.ipdatacloud.com/v2/query?ip=" + ip + "&key=" + secretKey
	client := http_util.HTTPClient{}
	data, err := client.Get(url)
	if err != nil {
		return err
	}
	var ipd IPLocation
	if err := json.Unmarshal(data, &ipd); err != nil {
		return err
	}
	if ipd.Code != 200 {
		return errors.New(fmt.Sprintf("获取出错 code:{%d} msg{%s}", ipd.Code, ipd.Msg))
	}
	*location = ipd.Data
	return nil
}

// 获取外网ip地址
func GetLocation(ip, key string) string {
	if ip == "127.0.0.1" || ip == "localhost" {
		return "内部IP"
	}
	url := "https://restapi.amap.com/v5/ip?ip=" + ip + "&type=4&key=" + key
	fmt.Println("url", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("restapi.amap.com failed:", err)
		return "未知位置"
	}
	defer resp.Body.Close()
	s, err := io.ReadAll(resp.Body)
	fmt.Println(string(s))

	m := make(map[string]string)

	err = json.Unmarshal(s, &m)
	if err != nil {
		fmt.Println("Umarshal failed:", err)
	}
	//if m["province"] == "" {
	//	return "未知位置"
	//}
	return m["country"] + "-" + m["province"] + "-" + m["city"] + "-" + m["district"] + "-" + m["isp"]
}

// 获取局域网ip地址
func GetLocalHost() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}

	}
	return ""
}
