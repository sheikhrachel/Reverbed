package main

import (
	"encoding/json"
	"os"
	"runtime"

	"github.com/sheikhrachel/reverbed/api_common/call"
	"github.com/sheikhrachel/reverbed/api_common/router"
	"github.com/sheikhrachel/reverbed/api_common/utils/errutil"
	"github.com/sheikhrachel/reverbed/api_common/utils/httpUtil"
	"github.com/sheikhrachel/reverbed/handlers"
)

func main() {
	// set the cpu cores config
	runtime.GOMAXPROCS(runtime.NumCPU())
	// setting up the router
	r := router.CreateRouterWithMiddleware()
	cc := router.InitCall()
	// setting up the routes
	// - this will set up the routes from
	//   the handlers package, within handlers/routes.go
	handlers.SetRoutes(r, cc)
	// setting up swagger
	r.Static("/swagger", "./swaggerui")
	// starting the router
	// - this will start accepting results on port 8080
	router.StartRouter(r, cc)
}

const weatherURL = "http://api.weatherapi.com/v1"

// getWeather returns the weather for Seattle, WA using the Weather API
func getWeather(cc call.Call) (weather WeatherAPIResponse) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return WeatherAPIResponse{}
	}
	body, err := httpUtil.GetRespBody(cc, apiKey, weatherURL+"/current.json?key="+apiKey+"&q=Seattle", nil)
	if errutil.HandleError(cc, err) {
		return WeatherAPIResponse{}
	}
	if err = json.Unmarshal(body, &weather); errutil.HandleError(cc, err) {
		return WeatherAPIResponse{}
	}
	return weather
}

type WeatherAPIResponse struct {
	Current WeatherAPIResponseCurrent `json:"current"`
}

type WeatherAPIResponseCurrent struct {
	TempF float64 `json:"temp_f"`
}
