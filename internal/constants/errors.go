package constants

import "fmt"

var ErrorDeviceNotFound = fmt.Errorf("device not found for given ID")
var ErrorInvalidIDFormat = fmt.Errorf("ID value must be an integer")
