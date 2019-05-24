package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormattedContact(t *testing.T) {
	c := Contato{
		Celular:     953200587,
		DDDCelular:  11,
		Telefone:    29435002,
		DDDTelefone: 11,
	}
	if "(11) 953200587" != c.GetFormattedCellPhone() {
		t.Fail()
	}

	if "(11) 29435002" != c.GetFormattedPhone() {
		t.Fail()
	}
}

func TestClassificacao(t *testing.T) {
	t.Run("Crianca", func(t *testing.T) {
		assert.Equal(t, "Criança", Membro{
			Pessoa: Pessoa{
				DtNascimento: time.Now(),
			},
		}.Classificacao())
	})
	t.Run("Adolescente", func(t *testing.T) {
		assert.Equal(t, "Adolescente", Membro{
			Pessoa: Pessoa{
				DtNascimento: time.Now().AddDate(-17, 0, 0),
			},
		}.Classificacao())
	})
	t.Run("Jovem", func(t *testing.T) {
		assert.Equal(t, "Jovem", Membro{
			Pessoa: Pessoa{
				DtNascimento: time.Now().AddDate(-29, 0, 0),
			},
		}.Classificacao())
	})
	t.Run("Adulto Solteiro", func(t *testing.T) {
		assert.Equal(t, "Adulto", Membro{
			Pessoa: Pessoa{
				DtNascimento: time.Now().AddDate(-33, 0, 0),
			},
		}.Classificacao())
	})
	t.Run("Adulto Casado", func(t *testing.T) {
		assert.Equal(t, "Adulto", Membro{
			Pessoa: Pessoa{
				DtNascimento: time.Now().AddDate(-25, 0, 0),
				DtCasamento:  time.Now(),
			},
		}.Classificacao())
	})
}

func TestFormattedAddress(t *testing.T) {
	address := Endereco{
		Logradouro: "Rua xicas",
		Bairro:     "Parque feliz",
		Numero:     2,
	}
	assert.Equal(t, "Rua xicas, 2 - Parque feliz", address.GetFormatted())
}
