package resources

type Device struct {
	ID           int    `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	Name         string `json:"name,omitempty"`
	Brand        string `json:"brand,omitempty"`
	State        State  `json:"state,omitempty"`
	CreationTime string `json:"creation_time,omitempty"`
}

type State string
