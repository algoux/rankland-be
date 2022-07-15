package utils

import "github.com/bwmarrin/snowflake"

var Generator *snowflake.Node

func InitGenerator() error {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return err
	}
	Generator = node

	return nil
}
