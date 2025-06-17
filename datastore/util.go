package datastore

type FilterData struct {
	Filters map[string]interface{} `json:"filters"`
}

func (obj *FilterData) Validate() error {
	return nil
}
