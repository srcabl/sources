package service

import (
	pb "github.com/srcabl/protos/shared"
)

// SourceDeterminer defines the behavior of a source determiner
type SourceDeterminer interface {
	DetermineSource(string) ([]*pb.SourceNode, error)
}

type sourceDeterminer struct {
	datarepo DataRepository
}

// NewSourceDeterminer news up a source deteminer
func NewSourceDeterminer(datarepo DataRepository) (SourceDeterminer, error) {
	return &sourceDeterminer{
		datarepo: datarepo,
	}, nil
}

// DetermineSource determines the source of a link
func (s *sourceDeterminer) DetermineSource(url string) ([]*pb.SourceNode, error) {
	//TODO: fillthis out
	return []*pb.SourceNode{
		{
			Source: &pb.Source{
				Name:         "Blahhh",
				Organization: "Booooo",
			},
		},
	}, nil
}
