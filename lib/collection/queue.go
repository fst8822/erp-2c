package collection

import "erp-2c/model"

type Queue struct {
	in  chan model.DeliveryDB
	out chan model.DeliveryDB
}

func NewQueue(capacity int) *Queue {
	return &Queue{
		in:  make(chan model.DeliveryDB, capacity),
		out: make(chan model.DeliveryDB, capacity),
	}
}

func (q *Queue) In() chan model.DeliveryDB {
	return q.in
}
func (q *Queue) Out() chan model.DeliveryDB {
	return q.out
}
func (q *Queue) CloseIn() {
	close(q.in)
}
func (q *Queue) CloseOut() {
	close(q.out)
}
