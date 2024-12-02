package data

import (
	"fmt"
	"github.com/goalm/kit/sys"
	"github.com/goalm/pGo/lexer"
	"github.com/goalm/pGo/token"
	"log"
	"strconv"
	"strings"
)

type Product struct {
	Idx        int
	ProdName   string
	Dimensions map[string]int
	Params     map[string]string
}

type ProductProperties struct {
	Idx        int               `csv:"IDX"`
	ProdName   string            `csv:"PROD_NAME"`
	Dimensions map[string]string `csv:"-"`
}

type ResultSpecs struct {
	Idx       int               `csv:"IDX"`
	ProdName  string            `csv:"PROD_NAME"`
	SpCode    string            `csv:"SP_CODE"`
	OtherData map[string]string `csv:"-"`
}

type VariableProperties struct {
	Idx        int               `csv:"IDX"`
	Formula    string            `csv:"FORMULA"`
	TimePoints string            `csv:"TIME_POINTS"`
	OtherData  map[string]string `csv:"-"`
}

type PathTrack struct {
	Idx           int               `csv:"IDX"`
	StartVariable string            `csv:"START_VARIABLE"`
	EndVariable   string            `csv:"END_VARIABLE"`
	OtherData     map[string]string `csv:"-"`
}

func ListFormulas(v *VariableProperties, p *ResultSpecs, prodDims map[string]map[string]string) []string {
	dimExpression := strings.Replace(v.OtherData["DIMENSIONS"], " ", "", -1)
	dimArrays := make([][]int, 0)
	formula := v.Formula
	formulas := make([]string, 0)

	if dimExpression != "" && dimExpression != "N/A" && dimExpression != "n/a" {
		dimensions := strings.Split(dimExpression, ",")
		for _, dim := range dimensions {
			value, ok := prodDims[p.ProdName][dim]
			if !ok {
				fmt.Printf("Dimension  %s not found in %s\n", dim, p.ProdName)
				continue
			}
			dimSlice, err := sys.ParseRange(value)
			if err != nil {
				fmt.Println(err)
			}
			dimArrays = append(dimArrays, dimSlice)
			log.Println(v.Formula, dimArrays)
		}

		// generate indices
		arrayList := sys.GenArryList(dimArrays)
		for _, idx := range arrayList {
			var vName string
			// parse formula
			l := lexer.New(formula)
			for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
				// add suffix to variable name
				if tok.Type == token.IDENT {
					vName += tok.Literal + "("
					for _, v := range idx[:len(idx)-1] {
						vName += strconv.Itoa(v) + ":"
					}
					vName += strconv.Itoa(idx[len(idx)-1]) + ")"
				} else {
					vName += tok.Literal
				}
			}
			formulas = append(formulas, vName)
		}
		return formulas
	} else {
		return []string{v.Formula}
	}
}
