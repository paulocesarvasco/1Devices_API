package constants

import "fmt"

var ErrorBrandNotFound = fmt.Errorf("any device found for given brand")
var ErrorStateNotFound = fmt.Errorf("any device found for given state")
var ErrorDeviceNotFound = fmt.Errorf("device not found for given ID")
var ErrorInvalidIDFormat = fmt.Errorf("ID value must be an integer")
var ErrorInvalidRequestParameter = fmt.Errorf("unexpected requested parameter")
var ErrorDeviceInUse = fmt.Errorf("requested device currently in use")
var ErrorMissedRequestIDParameter = fmt.Errorf("missed request device id")
var ErrorInvalidDeviceState = fmt.Errorf("unexpected given state")
