package watcher

/*
	Feature B:
		Counter struct now contains a HexValue field, in order to return the generated hex value.
*/

type Counter struct {
	Iteration int    `json:"iteration"`
	HexValue  string `json:"value"`
}

type CounterReset struct {
}
