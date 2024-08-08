package dynamodbhelper

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gofiber/fiber/v2/log"
	"reflect"
	"strings"
)

type UpdateMatcher struct {
	Table    string
	ID       string
	ChurchID string
	Values   map[string]types.AttributeValue
}

func (expected UpdateMatcher) Matches(r any) bool {
	received := r.(*dynamodb.UpdateItemInput)
	if *received.TableName != expected.Table {
		return false
	}
	if received.Key["id"].(*types.AttributeValueMemberS).Value != expected.ID {
		return false
	}
	if received.Key["church_id"].(*types.AttributeValueMemberS).Value != expected.ChurchID {
		return false
	}
	if received.UpdateExpression == nil {
		return false
	}
	updateQuery := *received.UpdateExpression
	if !reflect.DeepEqual(received.ExpressionAttributeValues, expected.Values) {
		return false
	}
	for key := range received.ExpressionAttributeValues {
		if !strings.Contains(updateQuery, key) {
			log.Errorf("Key %s not found in update query", key)
			return false
		}
	}

	return true
}

func (expected UpdateMatcher) String() string {
	return fmt.Sprintf("Expected ID: {%s}, ChurchID: {%s}, Table: {%s}, Values:{%v}", expected.ID, expected.ChurchID, expected.Table, expected.Values)
}
