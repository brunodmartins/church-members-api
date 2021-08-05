package file_test

import (
	"github.com/BrunoDM2943/church-members-api/internal/modules/report/file"
	"github.com/spf13/viper"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func init(){
	viper.Set("bundles.location", "../../../../bundles")
	viper.Set("pdf.font.path", "../../../../fonts/Arial.ttf")
}


func TestBuildFile(t *testing.T) {
	pdfBuilder := file.NewPDFBuilder()
	out, err := pdfBuilder.BuildFile("Test", BuildMembers(100))
	assert.False(t, utf8.Valid(out))
	assert.NotNil(t, out)
	assert.Nil(t, err)
}

func TestBuildFileErrorOnFont(t *testing.T) {
	viper.Set("pdf.font.path", ".")
	pdfBuilder := file.NewPDFBuilder()
	_, err := pdfBuilder.BuildFile("Test", BuildMembers(100))
	assert.NotNil(t, err)
}