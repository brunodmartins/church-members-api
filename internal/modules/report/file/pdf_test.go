package file_test

import (
	"testing"
	"unicode/utf8"

	"github.com/brunodmartins/church-members-api/internal/modules/report/file"
	"github.com/spf13/viper"

	"github.com/stretchr/testify/assert"
)

func init() {
	viper.Set("pdf.font.path", "../../../../resources/fonts/Arial.ttf")
}

func TestBuildFile(t *testing.T) {
	pdfBuilder := file.NewPDFBuilder()
	out, err := pdfBuilder.BuildFile("Test", buildChurch(), BuildMembers(100))
	assert.False(t, utf8.Valid(out))
	assert.NotNil(t, out)
	assert.Nil(t, err)
}

func TestBuildFileErrorOnFont(t *testing.T) {
	viper.Set("pdf.font.path", ".")
	pdfBuilder := file.NewPDFBuilder()
	_, err := pdfBuilder.BuildFile("Test", buildChurch(), BuildMembers(100))
	assert.NotNil(t, err)
}
