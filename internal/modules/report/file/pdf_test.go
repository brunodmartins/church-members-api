package file_test

import (
	"context"
	"testing"
	"unicode/utf8"

	"github.com/brunodmartins/church-members-api/internal/modules/report/file"
	"github.com/stretchr/testify/assert"
)

func TestBuildFile(t *testing.T) {
	pdfBuilder := file.NewPDFBuilder()
	out, err := pdfBuilder.BuildFile(context.Background(), "Test", buildChurch(), BuildMembers(100))
	assert.False(t, utf8.Valid(out))
	assert.NotNil(t, out)
	assert.Nil(t, err)
}
