package globalkey

/*
*
global constant key
*/
const (
	DefaultPageSize = 5
	MaxPageSize     = 10
)

// DelStateNo 软删除
var DelStateNo int64 = 0  //未删除
var DelStateYes int64 = 1 //已删除

// DateTimeFormatTplStandardDateTime 时间格式化模版
var DateTimeFormatTplStandardDateTime = "Y-m-d H:i:s"
var DateTimeFormatTplStandardDate = "Y-m-d"
var DateTimeFormatTplStandardTime = "H:i:s"
