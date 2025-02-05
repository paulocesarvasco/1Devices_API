package resources

type Device struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Brand        string `json:"brand"`
	State        string `json:"state"`
	CreationTime string `json:"creation_time"`
}
