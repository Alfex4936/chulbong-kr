package main

import (
	"encoding/json"
	"io"
	"math/rand"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/render"
)

var (
	baseMapData = []opts.MapData{
		{Name: "서울특별시", Value: float64(rand.Intn(150))},
		{Name: "부산광역시", Value: float64(rand.Intn(150))},
		{Name: "대구광역시", Value: float64(rand.Intn(150))},
		{Name: "인천광역시", Value: float64(rand.Intn(150))},
		{Name: "광주광역시", Value: float64(rand.Intn(150))},
		{Name: "대전광역시", Value: float64(rand.Intn(150))},
		{Name: "울산광역시", Value: float64(rand.Intn(150))},
		{Name: "세종특별자치시", Value: float64(rand.Intn(150))},
		{Name: "경기도", Value: float64(rand.Intn(150))},
		{Name: "강원도", Value: float64(rand.Intn(150))},
		{Name: "충청북도", Value: float64(rand.Intn(150))},
		{Name: "충청남도", Value: float64(rand.Intn(150))},
		{Name: "전라북도", Value: float64(rand.Intn(150))},
		{Name: "전라남도", Value: float64(rand.Intn(150))},
		{Name: "경상북도", Value: float64(rand.Intn(150))},
		{Name: "경상남도", Value: float64(rand.Intn(150))},
		{Name: "제주특별자치도", Value: float64(rand.Intn(150))},
	}

	seoulMapData = map[string]float64{
		"종로구":  float64(rand.Intn(150)),
		"중구":   float64(rand.Intn(150)),
		"용산구":  float64(rand.Intn(150)),
		"성동구":  float64(rand.Intn(150)),
		"광진구":  float64(rand.Intn(150)),
		"동대문구": float64(rand.Intn(150)),
		"중랑구":  float64(rand.Intn(150)),
		"성북구":  float64(rand.Intn(150)),
		"강북구":  float64(rand.Intn(150)),
		"도봉구":  float64(rand.Intn(150)),
		"노원구":  float64(rand.Intn(150)),
		"은평구":  float64(rand.Intn(150)),
		"서대문구": float64(rand.Intn(150)),
		"마포구":  float64(rand.Intn(150)),
		"양천구":  float64(rand.Intn(150)),
		"강서구":  float64(rand.Intn(150)),
		"구로구":  float64(rand.Intn(150)),
		"금천구":  float64(rand.Intn(150)),
		"영등포구": float64(rand.Intn(150)),
		"동작구":  float64(rand.Intn(150)),
		"관악구":  float64(rand.Intn(150)),
		"서초구":  float64(rand.Intn(150)),
		"강남구":  float64(rand.Intn(150)),
		"송파구":  float64(rand.Intn(150)),
		"강동구":  float64(rand.Intn(150)),
	}
)


func generateMapData(data map[string]float64) (items []opts.MapData) {
	items = make([]opts.MapData, 0)
	for k, v := range data {
		items = append(items, opts.MapData{Name: k, Value: v})
	}
	return
}

func mapBase() *charts.Map {
	mc := charts.NewMap()
	mc.RegisterMapType("south_korea")
	mc.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "basic map example"}),
	)

	mc.AddSeries("map", baseMapData)
	mc.AddJSFuncs(render.CustomizedJSAssets("south_korea.js"))
	return mc
}

func mapShowLabel() *charts.Map {
	mc := charts.NewMap()
	mc.RegisterMapType("south_korea")
	mc.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "show label"}),
	)

	mc.AddSeries("map", baseMapData).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show: true,
			}),
		)
	mc.AddJSFuncs(render.CustomizedJSAssets("south_korea.js"))
	return mc
}

func mapVisualMap() *charts.Map {
	mc := charts.NewMap()
	mc.RegisterMapType("south_korea")
	mc.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "VisualMap",
		}),
		charts.WithVisualMapOpts(opts.VisualMap{
			Calculable: true,
		}),
	)

	mc.AddSeries("map", baseMapData)
	mc.AddJSFuncs(render.CustomizedJSAssets("south_korea.js"))
	return mc
}

func mapRegion() *charts.Map {
	mc := charts.NewMap()
	mc.RegisterMapType("서울")
	mc.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Seoul province",
		}),
		charts.WithVisualMapOpts(opts.VisualMap{
			Calculable: true,
			InRange:    &opts.VisualMapInRange{Color: []string{"#50a3ba", "#eac736", "#d94e5d"}},
		}),
	)

	mc.AddSeries("map", generateMapData(seoulMapData))
	mc.AddJSFuncs(render.CustomizedJSAssets("south_korea.js"))
	return mc
}

func mapTheme() *charts.Map {
	mc := charts.NewMap()
	mc.RegisterMapType("south_korea")
	mc.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme: "macarons",
		}),
		charts.WithTitleOpts(opts.Title{
			Title: "Map-theme",
		}),
		charts.WithVisualMapOpts(opts.VisualMap{
			Calculable: true,
			Max:        150,
		}),
	)

	mc.AddSeries("map", baseMapData)
	mc.AddJSFuncs(render.CustomizedJSAssets("south_korea.js"))
	return mc
}

func main() {
	if err != nil {
		panic(err)
	}

	page := components.NewPage()
	page.AddCharts(
		mapBase(),
		mapShowLabel(),
		mapVisualMap(),
		mapRegion(),
		mapTheme(),
	)

	f, err := os.Create("map.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}
