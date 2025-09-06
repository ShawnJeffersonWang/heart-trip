package uniqueid

import (
	"fmt"
	"heart-trip/common/tool"
	"time"
)

// SnPrefix 生成sn单号
type SnPrefix string

const (
	SnPrefixHomestayOrder SnPrefix = "HSO" //民宿订单前缀 heart_trip_order/homestay_order
	SnPrefixThirdPayment  SnPrefix = "PMT" //第三方支付流水记录前缀 heart_trip_payment/third_payment
)

// GenSn 生成单号
func GenSn(snPrefix SnPrefix) string {
	return fmt.Sprintf("%s%s%s", snPrefix, time.Now().Format("20060102150405"), tool.Krand(8, tool.KC_RAND_KIND_NUM))
}
