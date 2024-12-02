package prop

import (
	"github.com/goalm/api"
	"github.com/goalm/lib/utils"
	"log"
	"math"
	"strconv"
	"strings"
)

const (
	zeroTol = 1e-20
)

type ProphetRecord struct {
	Name    string
	ResSpec string
	WsLoc   string
	RunNo   string
	Product string
	SpCode  string
	VarId   int
	VarName string
	Res     []float64
}

func ComposeEplHeader(startYr, startMth, stepSize, NumSteps, yearsOfMonthlyResults int) []string {
	header := []string{"!4", "ProdSp", "IdxCflow", "VarName"}
	for i := 0; i < NumSteps; i++ {
		start := utils.Date{startYr, startMth}
		Date := start.CalendarDate(i * stepSize)
		timePoint := Date.DateStr()
		if Date.Year >= startYr+yearsOfMonthlyResults {
			// skip non-December months
			if Date.Month != 12 {
				continue
			}
		}

		header = append(header, timePoint)
	}
	return header
}

func ComposeRecord(name, resSpec, wsLoc, runNo, pName, spCode, vName string, varId, startYr, startMth, stepSize, NumSteps, yearsOfMonthlyResults int) (r *ProphetRecord) {
	r = &ProphetRecord{}
	r.Name = name
	r.ResSpec = resSpec
	r.WsLoc = wsLoc
	r.RunNo = runNo
	r.Product = pName
	r.SpCode = spCode
	r.VarId = varId
	r.VarName = vName
	resSpecs := strings.Split(resSpec, "-")

	switch strings.TrimSpace(resSpecs[0]) {
	case "projResult":
		r.Res = ProphetProjResult(wsLoc, runNo, pName, spCode, vName, startYr, startMth, NumSteps, stepSize, yearsOfMonthlyResults)
	case "stoSummary":
		if len(resSpecs) != 2 {
			log.Fatalf("Invalid resSpec for session %v, please check your configuration", name)
		}
		stat := strings.TrimSpace(resSpecs[1])
		r.Res = ProphetStoSummary(wsLoc, runNo, pName, spCode, vName, stat, startYr, startMth, NumSteps, stepSize, yearsOfMonthlyResults)
	case "stoResult":
		if len(resSpecs) != 2 {
			log.Fatalf("Invalid resSpec for session %v, please check your configuration", name)
		}
		sim, err := strconv.Atoi(strings.TrimSpace(resSpecs[1]))
		if err != nil {
			log.Fatalf("Invalid resSpec for session %v, please check your configuration", name)
		}
		r.Res = ProphetStoResult(wsLoc, runNo, pName, spCode, vName, startYr, startMth, NumSteps, stepSize, sim, yearsOfMonthlyResults)
	}

	return
}

func ProphetProjResult(wsLoc, runNo, pName, sp, vName string, startYr, startMth, steps, stepSize, yearsOfMonthlyResults int) []float64 {
	ret := make([]float64, 0)
	for i := 0; i < steps; i++ {
		start := utils.Date{startYr, startMth}
		Date := start.CalendarDate(i * stepSize)
		timePoint := Date.DateStr()
		if Date.Year >= startYr+yearsOfMonthlyResults {
			// skip non-December months
			if Date.Month != 12 {
				continue
			}
			timePoint = strconv.Itoa(Date.Year)
		}

		tks, _ := utils.Parse(vName)
		res := 0.0
		opt := 1.0
		for _, v := range tks {
			if v.Tok == "+" {
				opt = 1.0
			} else if v.Tok == "-" {
				opt = -1.0
			} else {
				r := api.ReadProjResult(wsLoc, runNo, pName, sp, v.Tok, timePoint)
				if math.Abs(r) < zeroTol {
					r = 0
				}
				res += opt * r
			}
		}
		ret = append(ret, res)
	}
	return ret
}

func ProphetStoResult(wsLoc, runNo, pName, sp, vName string, startYr, startMth, steps, stepSize, sim, yearsOfMonthlyResults int) []float64 {
	ret := make([]float64, 0)
	for i := 0; i < steps; i++ {
		start := utils.Date{startYr, startMth}
		Date := start.CalendarDate(i * stepSize)
		timePoint := Date.DateStr()
		if Date.Year >= startYr+yearsOfMonthlyResults {
			// skip non-December months
			if Date.Month != 12 {
				continue
			}
			timePoint = strconv.Itoa(Date.Year)
		}

		tks, _ := utils.Parse(vName)
		res := 0.0
		opt := 1.0
		for _, v := range tks {
			if v.Tok == "+" {
				opt = 1.0
			} else if v.Tok == "-" {
				opt = -1.0
			} else {
				r := api.ReadStochasticResult(wsLoc, runNo, pName, sp, v.Tok, timePoint, sim)
				if math.Abs(r) < zeroTol {
					r = 0
				}
				res += opt * r
			}
		}
		ret = append(ret, res)
	}
	return ret
}

func ProphetStoSummary(wsLoc, runNo, pName, sp, vName, stat string, startYr, startMth, steps, stepSize, yearsOfMonthlyResults int) []float64 {
	ret := make([]float64, 0)
	for i := 0; i < steps; i++ {
		start := utils.Date{startYr, startMth}
		Date := start.CalendarDate(i * stepSize)
		timePoint := Date.DateStr()
		if Date.Year >= startYr+yearsOfMonthlyResults {
			// skip non-December months
			if Date.Month != 12 {
				continue
			}
			timePoint = strconv.Itoa(Date.Year)
		}

		tks, _ := utils.Parse(vName)
		res := 0.0
		opt := 1.0
		for _, v := range tks {
			if v.Tok == "+" {
				opt = 1.0
			} else if v.Tok == "-" {
				opt = -1.0
			} else {
				r := api.ReadStochasticSummary(wsLoc, runNo, pName, sp, v.Tok, timePoint, stat)
				if math.Abs(r) < zeroTol {
					r = 0
				}
				res += opt * r
			}
		}
		ret = append(ret, res)
	}
	return ret
}
