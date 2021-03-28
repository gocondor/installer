package cmd_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	. "github.com/gincoat/installer/gincoat/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/thanhpk/randstr"
)

func TestDownloadConfig(t *testing.T) {
	// Prepare
	newCmd := CmdNew{}
	var config Config
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		res, err := os.ReadFile("testdata/config.json")
		if err != nil {
			t.Fatal("error reading test data", err)
		}
		rw.Write(res)
	}))
	defer server.Close()

	// Execute
	newCmd.DownloadConfig(server.Client(), server.URL, &config)

	// Assert
	assert.Equal(t, "dummyVersion", config.InstallerReleasedVersion)
	assert.Equal(t, "dummyName", config.Releases["latest"].Name)
	assert.Equal(t, "dummyUrl", config.Releases["latest"].Url)
}
func TestDownloadGincoat(t *testing.T) {
	// Prepare
	newCmd := CmdNew{}
	fileName := "gincoat_temp_" + randstr.Hex(8) + ".tar.gz"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		res, err := os.ReadFile("testdata/gincoat.tar.gz")
		if err != nil {
			t.Fatal("error reading test data", err)
		}
		rw.Write(res)
	}))
	defer server.Close()

	// Execute
	filePath := newCmd.DownloadGincoat(server.Client(), server.URL, fileName)

	// Assert
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		t.Fatal("downloaded file not exist", filePath)
	}

	// Cleanup
	os.Remove(filePath)
}

func TestIsUpdatedRequired(t *testing.T) {
	// Prepare
	cn := CmdNew{}
	var config Config
	res, _ := os.ReadFile("testdata/config.json")
	json.Unmarshal(res, &config)

	// Execute
	check := cn.IsUpdatedRequired(config.InstallerReleasedVersion)

	// Assert
	if check != true { // supposed to be true
		t.Fatal("failed to assert check for update")
	}
}

func TestUnpack(t *testing.T) {
	// Prepare
	cn := CmdNew{}
	filepath := "./testdata/gincoat.tar.gz"
	folderName := "gincoat-0.3-alpha.3"
	destPath := os.TempDir()
	os.RemoveAll(destPath + "/" + folderName)

	// Execute
	cn.Unpack(filepath, destPath)

	// Assert
	_, err := os.Stat(destPath + "/" + folderName)
	if os.IsNotExist(err) {
		t.Fatal("failed to assert unpack")
	}

	files, err := ioutil.ReadDir(destPath + "/" + folderName)
	if len(files) <= 0 {
		t.Fatal("failed to assert unpack")
	}

	// remove the temp dir
	os.RemoveAll(destPath + "/" + folderName)
}
