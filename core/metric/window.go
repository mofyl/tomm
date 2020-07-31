package metric

type Bucket struct {
	Points []float64
	Count  int64
	// 这里是环状的，使用数组保存数据，然后 将每个 Bucket链接起来
	next *Bucket
}

func (b *Bucket) Append(val float64) {
	b.Points = append(b.Points, val)
	b.Count++
}

func (b *Bucket) Add(offset int, val float64) {
	b.Points[offset] += val
	b.Count++
}

func (b *Bucket) Reset() {
	b.Points = b.Points[:0]
	b.Count = 0
}

func (b *Bucket) Next() *Bucket {
	return b.next
}

type Window struct {
	/*
		这里的window是一个 时间轮 利用 Bucket中的next将所有的Bucket串起来
		将采样时间 分配到每个bucket中，每个bucket就是一个时间格子。
		然后 根据offset 计算出当前时间应该在那个格子中
	*/
	window []Bucket
	size   int
}

func NewWindow(size int) *Window {

	buckets := make([]Bucket, size)

	for i := range buckets {
		buckets[i] = Bucket{
			Points: make([]float64, 0),
		}
		next := i + 1
		if next == size {
			next = 0
		}

		buckets[i].next = &buckets[next]
	}

	return &Window{
		window: buckets,
		size:   size,
	}
}

func (w *Window) ResetWindow() {

	for i := range w.window {
		w.ReSetBucket(i)
	}

}

func (w *Window) ReSetBucket(offset int) {
	w.window[offset].Reset()
}

func (w *Window) ResetBuckets(offsets []int) {
	for _, offset := range offsets {
		w.ReSetBucket(offset)
	}
}

func (w *Window) Append(offset int, val float64) {
	w.window[offset].Append(val)
}

func (w *Window) Add(offset int, val float64) {
	// 这里若Count == 0 那么就用第一个位置, 不知道为什么
	if w.window[offset].Count == 0 {
		w.window[offset].Append(val)
		return
	}
	// 后面的值都会在0号位置累加  但是后面计算的时候 0号位置 和1号位置的值都会计算
	w.window[offset].Add(0, val)
}

func (w *Window) Bucket(offset int) Bucket {
	return w.window[offset]
}

func (w *Window) Size() int {
	return w.size
}

func (w *Window) Iterator(offset int, count int) Iterator {
	return Iterator{
		count:         count,
		iteratedCount: 0,
		cur:           &w.window[offset],
	}
}
