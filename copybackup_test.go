package copybackup

import (
	"os"
	"reflect"
	"testing"
)

func TestGetFileList(t *testing.T) {
	path, _ := os.Getwd()
	files, _ := GetFileList(path)
	expected := []string{"/Users/yukimemi/.ghq/bitbucket.org/yukimemi/go/copybackup/copybackup.go", "/Users/yukimemi/.ghq/bitbucket.org/yukimemi/go/copybackup/copybackup_test.go", "/Users/yukimemi/.ghq/bitbucket.org/yukimemi/go/copybackup/export_test.go"}
	if !reflect.DeepEqual(files, expected) {
		t.Errorf("expected = [%s] but files = [%s]", expected, files)
	}
}
