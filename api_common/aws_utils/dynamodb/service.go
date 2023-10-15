package dynamodb

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/sheikhrachel/reverbed/api_common/call"
)

type Client interface {
	PutItemInTable(cc call.Call, tableName string, item interface{}, retriesLeft int) (err error)
	UpdateItemInTable(cc call.Call, tableName, keyName, keyValue, updateExpression, conditionalExpression string, attributeNames map[string]*string, key, attributeValues map[string]*dynamodb.AttributeValue, retriesLeft int) (err error)
	GetItemFromTable(cc call.Call, tableName, keyName, keyValue string, key map[string]*dynamodb.AttributeValue, resultInterface interface{}, retriesLeft int) (err error)
}
