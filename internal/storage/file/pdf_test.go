package file

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/model"
	"github.com/BrunoDM2943/church-members-api/internal/infra/i18n"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBuildFile(t *testing.T){
	i18n.BootStrapI18N()

	dtNascimento, _ := time.Parse("2006/01/02", "2020/06/07")
	dtCasamento, _ := time.Parse("2006/01/02", "2019/09/14")
	data := []*model.Member{
		{
			Person: model.Person{
				FirstName:    "Test",
				LastName:     "test test",
				BirthDate:    dtNascimento,
				MarriageDate: dtCasamento,
				SpousesName:  "Test spuse",
				Contact: model.Contact{
					CellPhoneArea: 99,
					CellPhone:     1234567890,
					PhoneArea:     99,
					Phone:         12345678,
					Email:         "teste@test.com",
				},
				Address: model.Address{
					District: "9",
					City:     "Does not sleep",
					State:    "My-State",
					Address:  "XXXXX",
					Number:   9,
				},
			},
		},
	}

	pdfBuilder := NewPDFBuilder()
	out, err := pdfBuilder.BuildFile("Test", data)

	assert.NotNil(t, out)
	assert.Nil(t, err)
}