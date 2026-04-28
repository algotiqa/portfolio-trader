//=============================================================================
/*
Copyright © 2025 Andrea Carboni andrea.carboni71@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
//=============================================================================

package simulation

import (
	"encoding/base64"
	"log/slog"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/algotiqa/portfolio-trader/pkg/core"
	"github.com/algotiqa/portfolio-trader/pkg/db"
	"github.com/algotiqa/types"
	"github.com/go-analyze/charts"
	"golang.org/x/exp/stats"
)

//=============================================================================
//===
//=== Simulation process
//===
//=============================================================================

type Process struct {
	ts       *db.TradingSystem
	trades   *[]db.Trade
	req      *Request
	result   *Result
	stopping bool
}

//=============================================================================

const SimStatusIdle     = "idle"
const SimStatusWaiting  = "waiting"
const SimStatusRunning  = "running"
const SimStatusComplete = "complete"

//=============================================================================

func NewProcess(ts *db.TradingSystem, trades *[]db.Trade, req *Request) *Process {
	return &Process{
		ts    : ts,
		trades: trades,
		req   : req,
		result: &Result{
			Status: SimStatusWaiting,
		},
	}
}

//=============================================================================

func (p *Process) Start() {
	slog.Info("SimulationProcess: Starting", "id", p.ts.Id)

	p.result = NewResult(p.GetFirstTradeDate(), p.GetLastTradeDate(), p.req.Runs, p.ts)

	p.result.GrossAll = run(p.req, p.trades, db.TradeTypeAll, 0)
	p.result.Step++
	if !p.stopping {
		p.result.GrossLong = run(p.req, p.trades, db.TradeTypeLong, 0)
		p.result.Step++
		if !p.stopping {
			p.result.GrossShort = run(p.req, p.trades, db.TradeTypeShort, 0)
			p.result.Step++
			if !p.stopping {
				p.result.NetAll = run(p.req, p.trades, db.TradeTypeAll, p.ts.CostPerOperation)
				p.result.Step++
				if !p.stopping {
					p.result.NetLong = run(p.req, p.trades, db.TradeTypeLong, p.ts.CostPerOperation)
					p.result.Step++
					if !p.stopping {
						p.result.NetShort = run(p.req, p.trades, db.TradeTypeShort, p.ts.CostPerOperation)
						p.result.Step++
					}
				}
			}
		}
	}

	p.result.Status = SimStatusComplete
	p.result.EndTime = time.Now()
	slog.Info("SimulationProcess: Ended", "id", p.ts.Id)
}

//=============================================================================

func (p *Process) Stop() {
	p.stopping = true
}

//=============================================================================

func (p *Process) GetResult() *Result {
	return p.result
}

//=============================================================================

func (p *Process) GetFirstTradeDate() types.Date {
	t := *p.trades
	return types.ToDate(t[0].ExitDate)
}

//=============================================================================

func (p *Process) GetLastTradeDate() types.Date {
	t := *p.trades
	return types.ToDate(t[len(*p.trades)-1].ExitDate)
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

func run(req *Request, trades *[]db.Trade, tradeType string, costPerOper float64) *Details {
	returns := core.GetReturns(trades, tradeType, costPerOper)
	size := len(returns)
	if size == 0 {
		return &Details{}
	}

	risk,err := core.CalcRisk(returns, costPerOper)
	if err != nil {
		return &Details{}
	}

	list         := core.CalcRMultiple(returns, risk)
	origEquity   := core.BuildEquity(&list)
	_, origDD    := core.BuildDrawDown(origEquity)
	origReturn   := (*origEquity)[size-1]
	origAvgTrade := origReturn / float64(size)

	sampleSet, equityReturns, maxDrawdowns := buildSampleSet(list, req.Runs)
	sampleSet = addMeanAndStdDev(sampleSet, size)
	sampleSet = append(sampleSet, *origEquity)

	medianReturn := stats.Median(equityReturns)
	medianMaxDD  := stats.Median(maxDrawdowns)
	ddDistrib    := buildDDDistrib(maxDrawdowns)
	ddProbab     := buildDDProbab (ddDistrib)

	p, err := buildChart(sampleSet, req.Width, req.Height)
	if err != nil {
		panic(err)
	}

	buf, err := p.Bytes()
	if err != nil {
		panic(err)
	}

	return &Details{
		DetectedRisk       : risk,
		NumberOfTrades     : len(list),
		EquitiesImage      : base64.StdEncoding.EncodeToString(buf),
		MaxDrawdownDistr   : ddDistrib,
		MaxDrawdownProb    : ddProbab,
		EquityReturn       : core.Trunc2d(origReturn),
		EquityMaxDD        : core.Trunc2d(origDD),
		EquityReturnDDRatio: calcRatio(origReturn, origDD),
		EquityAverageTrade : core.Trunc2d(origAvgTrade),
		MedianReturn       : core.Trunc2d(medianReturn),
		MedianMaxDD        : core.Trunc2d(medianMaxDD),
		MedianReturnDDRatio: calcRatio(medianReturn, medianMaxDD),
		MedianAverageTrade : core.Trunc2d(medianReturn / float64(size)),
	}
}

//=============================================================================

func buildSampleSet(list []float64, runs int) ([][]float64, []float64, []float64) {
	size := len(list)

	var sampleSet [][]float64
	var equityReturns []float64
	var maxDrawdowns  []float64

	for i := 0; i < runs; i++ {
		var sample = make([]float64, size)

		for j := 0; j < size; j++ {
			value := list[rand.Intn(size)]
			sample[j] = value
		}

		equity := core.BuildEquity(&sample)
		sampleSet = append(sampleSet, *equity)

		_, maxDD := core.BuildDrawDown(equity)
		maxDrawdowns = append(maxDrawdowns, maxDD)

		er := (*equity)[len(*equity) -1]
		equityReturns = append(equityReturns, er)
	}

	return sampleSet, equityReturns, maxDrawdowns
}

//=============================================================================

func addMeanAndStdDev(sampleSet [][]float64, size int) [][]float64 {
	mean, stdDev := buildMeanAndStdDev(sampleSet, size)
	up, down := buildUpAndLowStdDev(mean, stdDev)

	sampleSet = append(sampleSet, make([]float64, size))
	sampleSet = append(sampleSet, mean)
	sampleSet = append(sampleSet, up)
	sampleSet = append(sampleSet, down)

	return sampleSet
}

//=============================================================================

func buildMeanAndStdDev(sampleSet [][]float64, size int) ([]float64, []float64) {
	var mean   []float64
	var stdDev []float64

	for i := 0; i < size; i++ {
		var serie []float64

		for j := 0; j < len(sampleSet); j++ {
			serie = append(serie, sampleSet[j][i])
		}

		m, s := stats.MeanAndStdDev(serie)

		mean   = append(mean, m)
		stdDev = append(stdDev, s)
	}

	return mean, stdDev
}

//=============================================================================

func buildUpAndLowStdDev(mean, stdDev []float64) ([]float64, []float64) {
	var upStdDev   []float64
	var downStdDev []float64

	for i, v := range mean {
		upStdDev   = append(upStdDev,   v+stdDev[i])
		downStdDev = append(downStdDev, v-stdDev[i])
	}

	return upStdDev, downStdDev
}

//=============================================================================

func buildDDDistrib(maxDrawdowns []float64) *Distribution {
	minv := core.CalcMin(maxDrawdowns)
	size := int(math.Trunc(math.Abs(minv))) + 1

	var xAxis []string
	for i := 1; i <= size; i++ {
		xAxis = append(xAxis, strconv.Itoa(-size+i)+"R")
	}

	yAxis := make([]float64, size)

	for _, value := range maxDrawdowns {
		index := size + int(math.Trunc(value)) - 1
		yAxis[index]++
	}

	return &Distribution{
		XAxis: xAxis,
		YAxis: yAxis,
	}
}

//=============================================================================

func buildDDProbab(d * Distribution) *Distribution {
	yAxis := make([]float64, len(d.YAxis))

	for i, value := range d.YAxis {
		yAxis[i] = value
		if i>0 {
			yAxis[i] += yAxis[i -1]
		}
	}

	maxVal := yAxis[len(d.YAxis)-1]

	for i, value := range yAxis {
		yAxis[i] = math.Round(value * 100 / maxVal)
	}

	return &Distribution{
		XAxis: d.XAxis,
		YAxis: yAxis,
	}
}

//=============================================================================
//=== Chart building
//=============================================================================

func buildChart(sampleSet [][]float64, width, height int) (*charts.Painter, error) {
	xAxis := calcXAxis(len(sampleSet[0]))

	opt := charts.NewLineChartOptionWithData(sampleSet)
	opt.XAxis.Title  = "Trades"
	opt.XAxis.Labels = xAxis
	opt.XAxis.LabelFontStyle.FontSize = 8
	opt.YAxis[0].Title = "Cumulative R multiples"
	opt.YAxis[0].LabelFontStyle.FontSize = 8
	opt.LineStrokeWidth = 1
	opt.Theme = opt.Theme.WithSeriesColors(buildColors(len(sampleSet)))

	p := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        width,
		Height:       height,
	})
	err := p.LineChart(opt)

	return p, err
}

//=============================================================================

func calcXAxis(size int) []string {
	var axis = make([]string, size)

	for i := 1; i <= size; i++ {
		axis[i-1] = strconv.Itoa(i)
	}

	return axis
}

//=============================================================================

func buildColors(size int) []charts.Color {
	var list []charts.Color

	for i := 0; i < size-5; i++ {
		list = append(list, charts.Color{192, 192, 192, 96})
	}

	list = append(list, charts.Color{128, 128, 128, 255})
	list = append(list, charts.Color{ 16,  16,  16, 255})
	list = append(list, charts.Color{ 80,  80,  80, 255})
	list = append(list, charts.Color{ 80,  80,  80, 255})
	list = append(list, charts.Color{  0, 100, 200, 255})

	return list
}

//=============================================================================

func calcRatio(ret, dd float64) float64 {
	if dd != 0 {
		return -core.Trunc2d(ret/dd)
	}

	return 100
}

//=============================================================================
