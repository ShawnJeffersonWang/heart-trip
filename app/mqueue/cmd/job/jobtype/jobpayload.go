package jobtype

import "heart-trip/app/order/model"

// DeferCloseHomestayOrderPayload defer close homestay order
type DeferCloseHomestayOrderPayload struct {
	Sn string
}

// PaySuccessNotifyUserPayload pay success notify ws
type PaySuccessNotifyUserPayload struct {
	Order *model.HomestayOrder
}
