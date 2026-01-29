package collection

import "erp-2c/model"

type Queue struct {
	data chan model.DeliverDomain
}

func NewQueue(size int64) *Queue {
	return &Queue{data: make(chan model.DeliverDomain, size)}
}

func (q *Queue) Enqueue(delivery model.DeliverDomain) {
	q.data <- delivery
}
func (q *Queue) Dequeue() model.DeliverDomain {
	return <-q.data
}
func (q *Queue) Peek() model.DeliverDomain {
	return <-q.data
}
func (q *Queue) IsEmpty() bool {
	return len(q.data) == 0
}
func (q *Queue) Size() int {
	return len(q.data)
}
