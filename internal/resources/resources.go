package resources

type Device struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Brand        string `json:"brand,omitempty"`
	State        string `json:"state,omitempty"`
	CreationTime string `json:"creation_time,omitempty"`
}
