package src

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// API KEY d8b90697a8eb8e570a3e526e813307c0

//{
//"id": 703448,
//"name": "Kiev",
//"country": "UA",
//"coord": {
//"lon": 30.516666,
//"lat": 50.433334
//}
//},

const URL = "http://api.openweathermap.org/data/2.5/forecast?id=703448&APPID=d8b90697a8eb8e570a3e526e813307c0"

type WeatherForecastData struct {
	gorm.Model
	HumanTime          string
	WeatherDescription string
	UnixTimestamp      int64
	Temperature        float32
	Pressure           float32
	Cloud              int
	Wind               float32
	WindDirection      int
	Sunset             int32
	Sunrize            int32
}

type WeatherForecastDataResponse struct {
	API    string                `json:"api"`
	Total  int                   `json:"total"`
	Entity []WeatherForecastData `json:"entity"`
}

func WeatherForecast(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":

		type incWeatherForecastData struct {
			HumanTime          string
			WeatherDescription string
			UnixTimestamp      int64
			Temperature        float32
			Pressure           float32
			Cloud              int
			Wind               float32
			WindDirection      int
			Sunset             int32
			Sunrize            int32
		}

		var incomingData incWeatherForecastData

		err := json.NewDecoder(r.Body).Decode(&incomingData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var wfd WeatherForecastData

		Db.Where("unix_timestamp = ?", incomingData.UnixTimestamp).Find(&wfd)
		if wfd.WeatherDescription == "" {
			Db.Create(&WeatherForecastData{
				HumanTime:          incomingData.HumanTime,
				WeatherDescription: incomingData.WeatherDescription,
				UnixTimestamp:      incomingData.UnixTimestamp,
				Temperature:        incomingData.Temperature,
				Pressure:           incomingData.Pressure,
				Cloud:              incomingData.Cloud,
				Wind:               incomingData.Wind,
				WindDirection:      incomingData.WindDirection,
				Sunset:             incomingData.Sunset,
				Sunrize:            incomingData.Sunrize,
			})
		} else {
			Db.Model(&wfd).Update(WeatherForecastData{
				HumanTime:          incomingData.HumanTime,
				WeatherDescription: incomingData.WeatherDescription,
				Temperature:        incomingData.Temperature,
				Pressure:           incomingData.Pressure,
				Cloud:              incomingData.Cloud,
				Wind:               incomingData.Wind,
				WindDirection:      incomingData.WindDirection,
			})
		}

	case "GET":

		//SD := r.URL.Query().Get("start_date")
		//ED := r.URL.Query().Get("end_date")

		log.Println("/weather_forecast GET answered\n\n")

		var Response WeatherForecastDataResponse

		Db.Find(&Response.Entity)

		Response.API = Version
		Response.Total = len(Response.Entity)

		addedRecordString, _ := json.Marshal(Response)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		n, _ := fmt.Fprintf(w, string(addedRecordString))
		fmt.Println(n)

		log.Println("/weather_forecast GET answered\n\n")

	}
}

func GetWeatherForecast(t time.Time) {
	fmt.Printf("%v: Hello, World!\n", t)

	type sys struct {
		Pod string `json:"pod"`
	}
	type coord struct {
		Lat float32 `json:"lat"`
		Lon float32 `json:"lon"`
	}
	type city struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		Coord    coord  `json:"coord"`
		Country  string `json:"country"`
		Timezone int    `json:"timezone"`
		Sunrise  int32  `json:"sunrise"`
		Sunset   int32  `json:"sunset"`
	}
	type wind struct {
		Speed float32 `json:"speed"`
		Deg   int     `json:"deg"`
	}
	type clouds struct {
		All int `json:"all"`
	}
	type weather struct {
		Id          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}
	type main struct {
		Temp       float32 `json:"temp"`
		Feels_like float32 `json:"feels_like"`
		Temp_min   float32 `json:"temp_min"`
		Temp_max   float32 `json:"temp_max"`
		Pressure   float32 `json:"pressure"`
		Sea_level  int     `json:"sea_level"`
		Grnd_level int     `json:"grnd_level"`
		Humidity   int     `json:"humidity"`
		Temp_kf    float32 `json:"temp_kf"`
	}
	type OneHour struct {
		Dt      int64     `json:"dt"`
		Main    main      `json:"main"`
		Weather []weather `json:"weather"`
		Clouds  clouds    `json:"clouds"`
		Wind    wind      `json:"wind"`
		Sys     sys       `json:"sys"`
		Dt_txt  string    `json:"dt_txt"`
	}
	type jsonPocket struct {
		Cod     string    `json:"cod"`
		Message int       `json:"message"`
		Cnt     int       `json:"cnt"`
		List    []OneHour `json:"list"`
		City    city      `json:"city"`
	}

	resp, err := http.Get(URL)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(body))

	var response jsonPocket
	jsonErr := json.Unmarshal(body, &response)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	fmt.Println("\n")
	fmt.Println(response)

	for ind, val := range response.List {
		fmt.Println(ind)
		fmt.Println(val)

		tm := time.Unix(val.Dt, 0)
		fmt.Println(tm)

		var wfd WeatherForecastData

		Db.Where("unix_timestamp = ?", val.Dt).Find(&wfd)
		if wfd.WeatherDescription == "" {
			Db.Create(&WeatherForecastData{
				HumanTime:          time.Unix(val.Dt, 0).Format("2006-01-02 15:04:05"),
				WeatherDescription: val.Weather[0].Description,
				UnixTimestamp:      val.Dt,
				Temperature:        val.Main.Temp - 273.15,
				Pressure:           val.Main.Pressure,
				Cloud:              val.Clouds.All,
				Wind:               val.Wind.Speed,
				WindDirection:      val.Wind.Deg,
				Sunset:             response.City.Sunset,
				Sunrize:            response.City.Sunrise,
			})
		} else {
			Db.Model(&wfd).Where("unix_timestamp = ?", val.Dt).Update(WeatherForecastData{
				HumanTime:          time.Unix(val.Dt, 0).Format("2006-01-02 15:04:05"),
				WeatherDescription: val.Weather[0].Description,
				//UnixTimestamp:      val.Dt,
				Temperature:   val.Main.Temp - 273.15,
				Pressure:      val.Main.Pressure,
				Cloud:         val.Clouds.All,
				Wind:          val.Wind.Speed,
				WindDirection: val.Wind.Deg,
				//Sunset:             response.City.Sunset,
				//Sunrize:            response.City.Sunrise,
			})
		}

	}

}

func DoEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}
