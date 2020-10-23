package benchmark

import (
	"github.com/olesho/classify/sequence"
	"github.com/olesho/classify/stream"
	"testing"
)

func BenchmarkSequenceFile(b *testing.B) {
	r := sequence.NewRootCluster()
	r.LoadFile("../rozetka.html")
	b.StartTimer()
	r.Batch()
	b.StopTimer()
}

func BenchmarkStreamFile(b *testing.B) {
	s := stream.NewStorage()
	s.LoadFile("../rozetka.html")
	b.StartTimer()
	s.RunAsync()
	b.StopTimer()
}

func BenchmarkSequenceMultipleFiles(b *testing.B) {
	r := sequence.NewRootCluster()
	r.LoadFile("../rozetka1.html")
	r.LoadFile("../rozetka2.html")
	b.StartTimer()
	r.Batch()
	b.StopTimer()
}

func BenchmarkStreamMultipleFiles(b *testing.B) {
	s := stream.NewStorage()
	s.LoadFile("../rozetka1.html")
	s.LoadFile("../rozetka2.html")
	b.StartTimer()
	s.RunAsync()
	b.StopTimer()
}
