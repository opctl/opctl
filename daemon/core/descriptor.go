package core

type nodeDescriptor struct {
	Id        string
	OpRef     string
	OpGraphId string
	Op        *opDescriptor
	Container *containerDescriptor
}

type opDescriptor struct{}

type containerDescriptor struct{}
