package errutil

import "errors"

var (
	ErrDynamoGetItem        = errors.New("could not get item from dynamo")
	ErrDynamoGetItemResults = errors.New("empty dynamo item")
	ErrDynamoUnmarshalItem  = errors.New("could not unmarshall item")
)
