package data

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
