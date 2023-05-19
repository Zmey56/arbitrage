package getlocaldata

type ParameterForWorking struct {
	PayTypes        []string `json:"payTypes"`
	Amount          float64  `json:"transAmount"`
	PublisherType   string   `json:"publisher_type"`
	MonthOrderCount int      `json:"monthOrderCount"`
	MonthFinishRate float64  `json:"monthFinishRate"`
}
