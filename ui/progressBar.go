package ui

import "github.com/cheggaaa/pb/v3"

type progressBar struct {
	progressBar *pb.ProgressBar
}

type nilProgressBar struct{}

type ProgressBar interface {
	Inc()
	Finish()
	SetCurrent(value int)
}

// A simple progress bar CLI implementation
func NewProgressBar(count int) ProgressBar {
	p := pb.StartNew(count)

	return progressBar{
		progressBar: p,
	}
}

func (p progressBar) Inc() {
	p.progressBar.Increment()
}

func (p progressBar) Finish() {
	p.progressBar.Finish()
}

func (p progressBar) SetCurrent(value int) {
	p.progressBar.SetCurrent(int64(value))
}

func NilProgressBar() ProgressBar {
	return nilProgressBar{}
}

func (n nilProgressBar) Inc() {}

func (n nilProgressBar) Finish() {}

func (n nilProgressBar) SetCurrent(value int) {}
