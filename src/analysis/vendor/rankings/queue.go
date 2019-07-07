package rankings

import "sort"

type Queue struct {
	keys []int
	nodes map[int]int
}

func (q *Queue) Len() int {
	return len(q.keys)
}

func (q *Queue) Swap(i, j int) {
	q.keys[i], q.keys[j] = q.keys[j], q.keys[i]
}

func (q *Queue) Less(i, j int) bool {
	a := q.keys[i]
	b := q.keys[j]
	return q.nodes[a] < q.nodes[b]
}

func (q *Queue) Set(id int, priority int) {
	if _, ok := q.nodes[id]; !ok {
		q.keys = append(q.keys, id)
	}

	q.nodes[id] = priority

	sort.Sort(q)
}

func (q *Queue) Next() (id int, priority int) {
	key, keys := q.keys[0], q.keys[1:]
	q.keys = keys

	priority = q.nodes[key]

	delete(q.nodes, key)

	return key, priority
}

func (q *Queue) IsEmpty() bool { return len(q.keys) == 0 }

func (q *Queue) Get(id int) (priority int, ok bool) {
	priority, ok = q.nodes[id]
	return
}

func NewQueue() *Queue {
	var q Queue
	q.nodes = make(map[int]int)
	return &q
}
