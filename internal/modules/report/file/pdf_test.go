package file_test

import (
	"context"
	"os"
	"path/filepath"
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

func TestBuildSingleMemberFile(t *testing.T) {
	pdfBuilder := file.NewPDFBuilder()
	out, err := pdfBuilder.BuildSingleMemberFile(context.Background(), "Member Profile", buildChurch(), BuildMembers(1)[0])
	assert.False(t, utf8.Valid(out))
	assert.NotNil(t, out)
	assert.Nil(t, err)
	assert.Greater(t, len(out), 100)

	path := filepath.Join(os.TempDir(), "member-profile-review.pdf")
	err = os.WriteFile(path, out, 0o600)
	assert.NoError(t, err)
	assert.FileExists(t, path)
		// keep a local path for manual review; remove the file after test if you prefer cleanup
}
