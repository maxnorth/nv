package e2e

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/maxnorth/nv/test"
	"gopkg.in/yaml.v3"
)

type TestSubject struct {
	Describe string
	Test     []TestCase
	File     string
}

type TestCase struct {
	It   string
	With TestCaseDetails
}

type TestCaseDetails struct {
	Cmd   string
	Pwd   string
	Files map[string]string
	Out   string
	Err   string
	Exit  int
}

func RunTestFile(t *testing.T, fileName string) {
	setPath(path.Join(test.RootDir(), "/dist"))
	testSubject := loadTestDef(t, path.Join(test.RootDir(), "test/e2e", fileName))
	for _, testCase := range testSubject.Test {
		t.Run(testCase.It, func(t *testing.T) {
			runTestCase(t, testSubject, &testCase)
		})
	}
}

func loadTestDef(t *testing.T, filePath string) *TestSubject {
	yamlData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read test yaml file: %s", filePath)
	}

	var testDef TestSubject
	err = yaml.Unmarshal(yamlData, &testDef)
	if err != nil {
		t.Fatalf("failed to unmarshall test yaml file: %v", err)
	}

	testDef.File = filePath

	return &testDef
}

func runTestCase(t *testing.T, testSubject *TestSubject, testCase *TestCase) {
	if testCase.With.Pwd != "" && testCase.With.Files != nil {
		t.Fatal("invalid test config, only one of with.pwd or with.files can be set, not both")
	} else if testCase.With.Pwd == "" {
		tmpDir := createTmpDir(t, testCase)
		testCase.With.Pwd = tmpDir
	}

	exitCode, outstr, errstr := runCommand(t, testCase.With.Cmd, testCase.With.Pwd)

	if exitCode != testCase.With.Exit {
		t.Fatalf("error: actual exit code '%d' does not match expected '%d'", exitCode, testCase.With.Exit)
	}

	if errstr != testCase.With.Err {
		t.Fatalf("error: actual stderr does not match expected:\n%s", test.GetColoredDiff(errstr, testCase.With.Err))
	}

	if outstr != testCase.With.Out {
		t.Fatalf("error: actual stdout does not match expected:\n%s", test.GetColoredDiff(outstr, testCase.With.Out))
	}
}

func createTmpDir(t *testing.T, testCase *TestCase) string {
	tmpDir := t.TempDir()
	for file, content := range testCase.With.Files {
		f, err := os.Create(tmpDir + "/" + file)
		if err != nil {
			t.Fatal(err)
		}
		_, err = f.WriteString(content)
		if err != nil {
			t.Fatal(err)
		}
	}

	return tmpDir
}

func runCommand(t *testing.T, cmdStr, pwd string) (exitCode int, outstr string, errstr string) {
	cmdFields := strings.Fields(cmdStr)
	cmdName, args := cmdFields[0], []string{}
	if len(cmdFields) > 1 {
		args = cmdFields[1:]
	}

	cmd := exec.Command(cmdName, args...)

	if strings.HasPrefix(pwd, "/") {
		cmd.Dir = pwd
	} else {
		cmd.Dir = path.Join(test.RootDir(), pwd)
	}

	outbuf := bytes.NewBufferString("")
	errbuf := bytes.NewBufferString("")
	cmd.Stdout = outbuf
	cmd.Stderr = errbuf

	_ = cmd.Run()

	outbyte, err := io.ReadAll(outbuf)
	if err != nil {
		t.Fatal(err)
	}
	errbyte, err := io.ReadAll(errbuf)
	if err != nil {
		t.Fatal(err)
	}

	// remove trailing whitespace on lines
	r := regexp.MustCompile(`[^\S\r\n]+\n`)
	outbyte = r.ReplaceAll(outbyte, []byte("\n"))
	errbyte = r.ReplaceAll(errbyte, []byte("\n"))

	return cmd.ProcessState.ExitCode(), string(outbyte), string(errbyte)
}

func setPath(nvPath string) {
	path := os.Getenv("PATH")
	if !strings.Contains(path, nvPath) {
		os.Setenv("PATH", path+":"+nvPath)
	}
}
