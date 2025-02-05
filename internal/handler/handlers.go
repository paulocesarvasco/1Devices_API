package handler

type Handler interface {
	DeleteDevice()
	RegisterDevice()
	SearchDevice()
	UpdateDevice()
}
