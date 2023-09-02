package main

type MsgQueue struct {
	data []string
}

func NewMsgQueue() MsgQueue {
	return MsgQueue{data: make([]string, 0)}
}

func (q *MsgQueue) Push(m string) {
	q.data = append(q.data, m)
}

func (q *MsgQueue) Pop() string {
	defer func() {
		if len(q.data) == 1 {
			q.data = []string{}
		} else {
			q.data = q.data[1:]
		}
	}()
	return q.data[0]
}
