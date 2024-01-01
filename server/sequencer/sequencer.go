package sequencer

import (
	"github.com/bwmarrin/snowflake"
)

func init() {
	g, err := snowflake.NewNode(1)
	if err != nil {
		panic(err.Error())
	}
	defaultIdGenerator = &GenericIdGenerator{
		node: g,
	}
}

var defaultIdGenerator IdGenerator

type IdGenerator interface {
	GenerateId() string
}

type GenericIdGenerator struct {
	node *snowflake.Node
}

func (ig *GenericIdGenerator) GenerateId() string {
	return ig.node.Generate().String()
}

func DefaultIdGenerator() IdGenerator {
	return defaultIdGenerator
}
