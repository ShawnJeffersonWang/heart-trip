package config

// KqConf 必须要自己写一个KqConf, 不能直接用kq.KqConf, 有出现有的字段没有set
type KqConf struct {
	Brokers []string
	Topic   string
}
