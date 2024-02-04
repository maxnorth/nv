package e2e_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
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
	Cmd  string
	Dir  string
	Out  string
	Err  string
	Exit int
}

func Test_E2E_Help(t *testing.T) {
	runTestSubject(t, "help.yml")
}

func Test_E2E_Print(t *testing.T) {
	runTestSubject(t, "print.yml")
}

func Test_E2E_Run(t *testing.T) {
	runTestSubject(t, "run.yml")
}

func runTestSubject(t *testing.T, fileName string) {
	setPath(path.Join(test.RootDir(), "/dist"))
	testSubject := loadTestDef(t, path.Join(test.RootDir(), "test/e2e", fileName))
	for _, testCase := range testSubject.Test {
		t.Run(testCase.It, func(t *testing.T) {
			runTestCase(t, testSubject, &testCase)
		})
	}
}

func getTestFiles(t *testing.T) []string {
	globTarget := path.Join(test.RootDir(), "test/e2e") + "/*.yml"
	yamlPaths, err := filepath.Glob(globTarget)
	if err != nil {
		fmt.Printf("failed to glob test yaml files at path: %s", globTarget)
		t.FailNow()
	}

	return yamlPaths
}

func loadTestDef(t *testing.T, filePath string) *TestSubject {
	yamlData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("failed to read test yaml file: %s", filePath)
		t.FailNow()
	}

	var testDef TestSubject
	err = yaml.Unmarshal(yamlData, &testDef)
	if err != nil {
		fmt.Printf("failed to unmarshall test yaml file: %v", err)
		t.FailNow()
	}

	testDef.File = filePath

	return &testDef
}

func runTestCase(t *testing.T, testSubject *TestSubject, testCase *TestCase) {
	if testCase.With.Dir == "" {
		testCase.With.Dir = "test/contexts/standard"
	}

	exitCode, outstr, errstr := runCommand(t, testCase.With.Cmd, testCase.With.Dir)

	if exitCode != testCase.With.Exit {
		fmt.Println("error: exit code does not match")
		t.FailNow()
	}

	if errstr != "" && testCase.With.Err == "" {
		fmt.Println("error: command unexpectedly failed")
		t.FailNow()
	}

	if errstr != testCase.With.Err {
		fmt.Println("error: stderr does not match")
		t.FailNow()
	}

	if outstr != testCase.With.Out {
		fmt.Println("error: stdout does not match")
		t.FailNow()
	}
}

func runCommand(t *testing.T, cmdStr, dir string) (exitCode int, outstr string, errstr string) {
	cmdFields := strings.Fields(cmdStr)
	cmdName, args := cmdFields[0], []string{}
	if len(cmdFields) > 1 {
		args = cmdFields[1:]
	}

	cmd := exec.Command(cmdName, args...)

	cmd.Dir = path.Join(test.RootDir(), dir)

	outbuf := bytes.NewBufferString("")
	errbuf := bytes.NewBufferString("")
	cmd.Stdout = outbuf
	cmd.Stderr = errbuf

	err := cmd.Run()

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
	if strings.Index(path, nvPath) == -1 {
		os.Setenv("PATH", path+":"+nvPath)
	}
}
