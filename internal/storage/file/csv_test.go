package file

import (
	"github.com/BrunoDM2943/church-members-api/internal/constants/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTransformCSVToData(t *testing.T) {
	t1, _ := time.Parse("02/01/2006", "07/06/2020")
	t2, _ := time.Parse("02/01/2006", "22/03/2020")
	data := []Data{
		{
			Value: 	model.Member{
				Person: model.Person{
					FirstName: "Teste",
					LastName:  "Teste",
					BirthDate: t1,
				},
			},
		},
		{
			Value: 	model.Member{
				Person: model.Person{
					FirstName: "Teste 2",
					LastName:  "Teste 2",
					BirthDate: t2,
				},
			},
		},
	}
	csvOut := TransformToCSVData(data, []string{"Name", "Date"}, func(row Data) []string {
		member := row.Value.(model.Member)
		return []string{
			member.Person.GetFullName(),
			member.Person.BirthDate.Format("02/01"),
		}
	})
	assert.Equal(t, 3, len(csvOut))
	assert.Equal(t, []string{"Name", "Date"}, csvOut[0])
	assert.Equal(t, []string{"Teste 2 Teste 2", "22/03"}, csvOut[2])
	assert.Equal(t, []string{"Teste Teste", "07/06"}, csvOut[1])
}