package metric

type Iterator struct {
	count         int
	iteratedCount int
	cur           *Bucket
}

func (i *Iterator) Next() bool {
	return i.count != i.iteratedCount
}

func (i *Iterator) Bucket() Bucket {

	if !i.Next() {
		panic("")
	}

	bucket := *i.cur
	i.iteratedCount++
	i.cur = i.cur.Next()
	return bucket

}
