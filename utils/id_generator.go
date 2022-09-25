package utils

import (
	"github.com/bwmarrin/snowflake"
	"github.com/sirupsen/logrus"
)

var Generator *snowflake.Node

func init() {
	node, err := snowflake.NewNode(1)
	if err != nil {
		logrus.WithError(err).Fatalf("init snowflake id generator failed")
	}
	Generator = node
}
