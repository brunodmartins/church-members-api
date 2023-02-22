package graphql

import (
	"github.com/brunodmartins/church-members-api/internal/constants/domain"
	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolvers(t *testing.T) {
	member := buildMember(domain.NewID())
	t.Run("idResolver", func(t *testing.T) {
		value, err := idResolver(buildParams(member))
		assert.Nil(t, err)
		assert.Equal(t, member.ID, value)
	})
	t.Run("classificationResolver", func(t *testing.T) {
		value, err := classificationResolver(buildParams(member))
		assert.Nil(t, err)
		assert.Equal(t, member.Classification().String(), value)
	})
	t.Run("fullNameResolver", func(t *testing.T) {
		value, err := fullNameResolver(buildParams(member.Person))
		assert.Nil(t, err)
		assert.Equal(t, member.Person.GetFullName(), value)
	})
	t.Run("ageResolver", func(t *testing.T) {
		value, err := ageResolver(buildParams(member.Person))
		assert.Nil(t, err)
		assert.Equal(t, 0, value)
	})
	t.Run("marriageDateResolver", func(t *testing.T) {
		value, err := marriageDateResolver(buildParams(member.Person))
		assert.Nil(t, err)
		assert.Equal(t, member.Person.MarriageDate, value)
	})
	t.Run("marriageDateResolver - nil", func(t *testing.T) {
		value, err := marriageDateResolver(buildParams(domain.Person{}))
		assert.Nil(t, err)
		assert.Nil(t, value)
	})
	t.Run("cellPhoneResolver", func(t *testing.T) {
		value, err := cellPhoneResolver(buildParams(member.Person.Contact))
		assert.Nil(t, err)
		assert.Equal(t, member.Person.Contact.GetFormattedCellPhone(), value)
	})
	t.Run("phoneResolver", func(t *testing.T) {
		value, err := phoneResolver(buildParams(member.Person.Contact))
		assert.Nil(t, err)
		assert.Equal(t, member.Person.Contact.GetFormattedPhone(), value)
	})
	t.Run("fullAddressResolver", func(t *testing.T) {
		value, err := fullAddressResolver(buildParams(member.Person.Address))
		assert.Nil(t, err)
		assert.Equal(t, member.Person.Address.String(), value)
	})
}

func buildParams(source interface{}) graphql.ResolveParams {
	return graphql.ResolveParams{
		Source: source,
	}
}
