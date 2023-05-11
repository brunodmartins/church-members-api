package file_test

import (
	"context"
	"io/fs"
	"os"
	"testing"
	"unicode/utf8"

	"github.com/brunodmartins/church-members-api/internal/modules/report/file"
	"github.com/stretchr/testify/assert"
)

func TestBuildFile(t *testing.T) {
	pdfBuilder := file.NewPDFBuilder()
	members := BuildMembers(100)
	members = append(members, buildInactiveMember())
	out, err := pdfBuilder.BuildFile(context.Background(), "Test", buildChurch(), members)
	assert.False(t, utf8.Valid(out))
	assert.NotNil(t, out)
	assert.Nil(t, err)
	os.WriteFile("test.pdf", out, fs.ModeTemporary)
}
