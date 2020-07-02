package ecode

const (
	mm_fail = "mm fail"
)

var (
	// -9100 ~ -9200 为mm预留
	MMFail = NewErrWithMsg(mm_fail, addCode(-9000))
)
