package prop

import (
	"github.com/goalm/kit/read"
	"log"
	"strconv"
	"strings"
)

type Session struct {
	ID       int
	Name     string
	ExtrConf struct {
		StartYr               int `mapstructure:"startYr"`
		StartMth              int `mapstructure:"startMth"`
		NumSteps              int `mapstructure:"numSteps"`
		StepSizeMths          int `mapstructure:"stepSizeMths"`
		YearsOfMonthlyResults int `mapstructure:"yearsOfMonthlyResults"`
	} `mapstructure:"extrConf"`
	ModelResults []*ModelResult `mapstructure:"modelResults"`
	Comparisons  []*Comparison  `mapstructure:"comparisons"`
}

type ModelResult struct {
	Name     string `mapstructure:"name"`
	WsLoc    string `mapstructure:"wsLoc"`
	RunNo    string `mapstructure:"runNo"`
	ResSpec  string `mapstructure:"resSpec"`
	ProdConf string `mapstructure:"prodConf"`
	VarList  string `mapstructure:"varList"`
	ResList  string `mapstructure:"resList"`
}

type Comparison struct {
	Name     string `mapstructure:"name"`
	LeftRun  string `mapstructure:"leftRun"`
	RightRun string `mapstructure:"rightRun"`
}

func ValidationSessionInputs(s *Session) {
	for _, run := range s.ModelResults {
		// validate
		if run.ResSpec == "" {
			log.Fatalf("ResSpec is empty for session %v, run %v", s.Name, run.Name)
		}
		if run.WsLoc == "" {
			log.Fatalf("Workspace location is empty for session %v, run %v", s.Name, run.Name)
		}
		resSpec := strings.Split(run.ResSpec, "-")
		if len(resSpec) == 0 {
			log.Fatalf("Invalid resSpec for session %v, run %v", s.Name, run.Name)
		}
		switch strings.TrimSpace(resSpec[0]) {
		case "projResult":
			if len(resSpec) != 1 {
				log.Fatalf("Invalid resSpec for session %v, run %v", s.Name, run.Name)
			}
		case "stoSummary":
			if len(resSpec) != 2 {
				log.Fatalf("Invalid resSpec for session %v, run %v", s.Name, run.Name)
			}
			stat := strings.TrimSpace(resSpec[1])
			if stat != "MEAN_VALUE" &&
				stat != "STANDARD_DEV" &&
				stat != "MIN_VALUE" &&
				stat != "MIN_SIM" &&
				stat != "MAX_VALUE" &&
				stat != "MAX_SIM" &&
				stat != "MEDIAN_VALUE" &&
				stat != "MEDIAN_SIM" {
				log.Fatalf("Invalid stat %v for stoch summary session %v, run %v", stat, s.Name, run.Name)
			}
		case "stoResult":
			if len(resSpec) != 2 {
				log.Fatalf("Invalid resSpec for session %v, run %v", s.Name, run.Name)
			}
			sim := strings.TrimSpace(resSpec[1])
			// sim is a number
			if _, err := strconv.Atoi(sim); err != nil {
				log.Fatalf("Invalid resSpec for session %v, run %v", s.Name, run.Name)
			}
		}

		// validate resList
		if run.ResList == "" {
			log.Fatalf("ResList is empty for session %v, run %v", s.Name, run.Name)
		} else {
			prodDims := read.CsvToProductDimensions(run.ProdConf)
			resList := read.CsvToResultList(run.ResList)
			varList := read.CsvToVariableList(run.VarList)
			for _, res := range resList {
				if _, ok := prodDims[res.ProdName]; !ok {
					log.Fatalf("Product %v in resList not found in prodConf for session %v, run %v", res.ProdName, s.Name, run.Name)
				}
				for _, v := range varList {
					dimExpression := strings.Replace(v.OtherData["DIMENSIONS"], " ", "", -1)
					if dimExpression != "" && dimExpression != "N/A" && dimExpression != "n/a" {
						dimensions := strings.Split(dimExpression, ",")
						for _, dim := range dimensions {
							_, ok := prodDims[res.ProdName][dim]
							if !ok {
								log.Fatalf("Error: session %v, run %v, dimension %v not found in product %v", s.Name, run.Name, dim, res.ProdName)
							}
						}
					}
				}
			}
		}
	}
}
