package limit

//
//var num int32
//
//func Get() bool {
//	if atomic.LoadInt32(&num) == 0 {
//		return false
//	}
//
//	atomic.AddInt32(&num, -1)
//
//	log.Debug("Get %d", num)
//	return true
//}
//
//func put(num1 int32, ctx context.Context) {
//
//	ticker := time.NewTicker(5 * time.Second)
//
//	for {
//		select {
//		case <-ctx.Done():
//			ticker.Stop()
//			return
//		case <-ticker.C:
//			atomic.AddInt32(&num, num1)
//			log.Debug("put %d", num)
//		}
//
//	}
//
//}
