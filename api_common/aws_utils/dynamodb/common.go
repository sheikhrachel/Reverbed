package dynamodb

import (
	"errors"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	MaxDynamoRetries  = 3
	MaxBatchGetSize   = 100
	MaxBatchWriteSize = 25
)

func TableExists(tableName *string, client Client) bool {
	describeTableInput := &dynamodb.DescribeTableInput{
		TableName: tableName,
	}
	_, err := client.DescribeTable(describeTableInput)
	if err == nil {
		return true
	}
	var resourceNotFoundException *dynamodb.ResourceNotFoundException
	ok := errors.As(err, &resourceNotFoundException)
	return !ok
}

func MarshallIds(idKey string, ids []*string) []map[string]*dynamodb.AttributeValue {
	marshalledIds := make([]map[string]*dynamodb.AttributeValue, 0)
	for _, id := range ids {
		attrVal := &dynamodb.AttributeValue{
			S: id,
		}
		marshalledIds = append(marshalledIds, map[string]*dynamodb.AttributeValue{
			idKey: attrVal,
		})
	}
	return marshalledIds
}
