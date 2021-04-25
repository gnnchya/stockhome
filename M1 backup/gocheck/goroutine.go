package gocheck

func GoAmount(cnt int, ch chan<- int) {
	cnt++
	ch <- cnt
}
