package graphql

import (
	"errors"
	"github.com/BrunoDM2943/church-members-api/entity"
	"github.com/BrunoDM2943/church-members-api/member"
	"github.com/graphql-go/graphql"
)

type memberResolver struct {
	service member.Service
}

func newMemberResolver(service member.Service) memberResolver {
	return memberResolver{
		service,
	}
}

func filterPerson(pessoas []*entity.Membro, f func(*entity.Membro) bool) []entity.Membro {
	vsf := make([]entity.Membro, 0)
	for _, v := range pessoas {
		if f(v) {
			vsf = append(vsf, *v)
		}
	}
	return vsf

}

func getFilterAsBool(args map[string]interface{}, key string) (bool, error) {
	if args[key] != nil {
		return args[key].(bool), nil
	} else {
		return false, errors.New("Not Found")
	}
}

func getFilterAsString(args map[string]interface{}, key string) (string, error) {
	if args[key] != nil {
		return args[key].(string), nil
	} else {
		return "", errors.New("Not Found")
	}
}

func (this memberResolver) memberResolver(params graphql.ResolveParams) (interface{}, error) {
	members, _ := this.service.FindAll()
	return filterPerson(members, func(m *entity.Membro) bool {
		var valActive = true
		var valSex = true

		if val, err := getFilterAsBool(params.Args, "active"); err == nil {
			valActive = m.Active == val
		}
		if val, err := getFilterAsString(params.Args, "sexo"); err == nil {
			valSex = m.Pessoa.Sexo == val
		}
		return valSex && valActive
	}), nil
}
