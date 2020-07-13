package limit

import (
	"time"
)

type RouterData struct {
	RTT    int32
	Method string
	Router string
}

type Limit struct {
	RTTOnload   float32 // 没有延迟的RTT 时间
	RTTActual   float32 // 请求的实际RTT 时间
	curRtt      int64   // 每秒对 各个接口发来的Rtt做一次采样
	longtermRtt int64
	TokenNums   int64 // 令牌个数
	routerData  chan RouterData
	queueSize   int32 // 允许排队的数量
}

var limit *Limit

func init() {
	limit = newLimit()
}

func newLimit() *Limit {
	l := Limit{
		RTTOnload:   3, // 这里默认是3s
		RTTActual:   0,
		curRtt:      0,
		longtermRtt: 0,
		TokenNums:   400,
		routerData:  make(chan RouterData),
		queueSize:   200,
	}
	return &l
}

func (l *Limit) limitJob() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case limitData, ok := <-l.routerData:
			if !ok {
				return
			}
			l.RTTActual = float32(limitData.RTT)
		case <-ticker.C:

		}
	}
}

func (l *Limit) Allow() bool {

	//gradientRatio := max(1.0, min(2.0, l.RTTActual/l.RTTOnload))
	return false
}

func max(num1, num2 float32) float32 {
	if num1 > num2 {
		return num1
	}

	return num2
}

func min(num1, num2 float32) float32 {

	if num1 < num2 {
		return num1
	}

	return num2
}
