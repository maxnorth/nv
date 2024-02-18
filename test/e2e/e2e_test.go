package e2e_test

import (
	"testing"

	"github.com/maxnorth/nv/test/e2e"
)

func Test_E2E_EnvLoader(t *testing.T) {
	e2e.RunTestFile(t, "config-dotenv.yml")
}

func Test_E2E_NvConfig(t *testing.T) {
	e2e.RunTestFile(t, "config-nv.yml")
}

func Test_E2E_Help(t *testing.T) {
	e2e.RunTestFile(t, "cmd-help.yml")
}

func Test_E2E_Print(t *testing.T) {
	e2e.RunTestFile(t, "cmd-print.yml")
}

func Test_E2E_Run(t *testing.T) {
	e2e.RunTestFile(t, "cmd-run.yml")
}

func Test_E2E_Resolver(t *testing.T) {
	e2e.RunTestFile(t, "behavior-resolver.yml")
}
