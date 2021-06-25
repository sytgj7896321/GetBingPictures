package channel

func CreateWorker(id int, f func(id int, ch chan int)) chan<- int {
	ch := make(chan int)
	go f(id, ch)
	return ch
}
