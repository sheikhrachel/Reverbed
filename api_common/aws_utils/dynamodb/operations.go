package dynamodb

import (
	"errors"
	"math"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/sheikhrachel/reverbed/api_common/call"
	"github.com/sheikhrachel/reverbed/api_common/utils/errutil"
)

// PutItemInTable handles adding an item to dynamo
func (s *Service) PutItemInTable(cc call.Call, tableName string, item interface{}, retriesLeft int) (err error) {
	av, err := dynamodbattribute.MarshalMap(item)
	if errutil.HandleError(cc, err) {
		return err
	}
	if _, err = s.Svc.PutItem(&dynamodb.PutItemInput{Item: av, TableName: aws.String(tableName)}); errutil.HandleError(cc, err) {
		if retriesLeft > 0 {
			// exponential backoff - 10ms, 20ms, 40ms (+ random sub-millisecond jitter)
			time.Sleep(time.Duration(math.Pow(2, float64(retriesLeft))*10+rand.Float64()) * time.Millisecond)
			return s.PutItemInTable(cc, tableName, item, retriesLeft-1)
		}
		cc.InfoF("Retries exceeded for error putting item in table: %+v, item: %+v, err: %+v", tableName, item, err)
		return err
	}
	return nil
}

// UpdateItemInTable handles adding an item to dynamo
func (s *Service) UpdateItemInTable(
	cc call.Call,
	tableName, keyName, keyValue, updateExpression, conditionalExpression string,
	attributeNames map[string]*string,
	key, attributeValues map[string]*dynamodb.AttributeValue,
	retriesLeft int,
) (err error) {
	updateInput := &dynamodb.UpdateItemInput{TableName: aws.String(tableName), Key: key, ExpressionAttributeNames: attributeNames, ExpressionAttributeValues: attributeValues, UpdateExpression: aws.String(updateExpression)}
	if conditionalExpression != "" {
		updateInput.ConditionExpression = aws.String(conditionalExpression)
	}
	if _, err = s.Svc.UpdateItem(updateInput); errutil.HandleError(cc, err) {
		if retriesLeft > 0 {
			// exponential backoff - 10ms, 20ms, 40ms (+ random sub-millisecond jitter)
			time.Sleep(time.Duration(math.Pow(2, float64(retriesLeft))*10+rand.Float64()) * time.Millisecond)
			return s.UpdateItemInTable(cc, tableName, keyName, keyValue, updateExpression, conditionalExpression, attributeNames, key, attributeValues, retriesLeft-1)
		}
		cc.InfoF("Retries exceeded for error updating item in table: %+v, err: %+v", updateInput, err)
		return err
	}
	return nil
}

// GetItemFromTable handles adding an item to dynamo
func (s *Service) GetItemFromTable(
	cc call.Call,
	tableName, keyName, keyValue string,
	key map[string]*dynamodb.AttributeValue,
	resultInterface interface{},
	retriesLeft int,
) (err error) {
	if key == nil {
		err = errors.New("key is nil")
		cc.InfoF("Error getting item from table: %+v, err: %+v", tableName, err)
		return err
	}
	result, err := s.Svc.GetItem(&dynamodb.GetItemInput{TableName: aws.String(tableName), Key: key})
	if errutil.HandleError(cc, err) {
		if retriesLeft > 0 {
			// exponential backoff - 10ms, 20ms, 40ms (+ random sub-millisecond jitter)
			time.Sleep(time.Duration(math.Pow(2, float64(retriesLeft))*10+rand.Float64()) * time.Millisecond)
			return s.GetItemFromTable(cc, tableName, keyName, keyValue, key, resultInterface, retriesLeft-1)
		}
		cc.InfoF("Retries exceeded for error getting item from table: %+v, err: %+v", keyValue, err)
		return errutil.ErrDynamoGetItem
	}
	if result.Item == nil {
		return errutil.ErrDynamoGetItemResults
	}
	if err = dynamodbattribute.UnmarshalMap(result.Item, &resultInterface); errutil.HandleError(cc, err) {
		return errutil.ErrDynamoUnmarshalItem
	}
	return nil
}
