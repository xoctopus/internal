// Package internal defines makefile generator
// +genx:doc
package internal

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xoctopus/x/misc/defers"
	"github.com/xoctopus/x/misc/must"
	"github.com/xoctopus/x/slicex"
	"github.com/xoctopus/x/stringsx"
)

// Makefile go project Makefile generator
type Makefile struct {
	// TestIgnore patterns unit testing and coverage ignores
	// .pb.go _mock.go _genx_ testing.go example/ will be ignored default
	TestIgnore []string `json:"test_ignore" cmd:""`
	// FormatIgnore patterns code formating ignores
	// .git/ .xgo/ *.pb.go *_genx_* will be ignored default
	FormatIgnore []string `json:"format_ignore" cmd:""`
	// Env global env variables
	Env []string `json:"env" cmd:""`
	// HackTest hack test entry
	HackTest bool `json:"hack_test" cmd:",default=false"`
	// Depend
	Depends Depends `json:"depends" cmd:"depends"`

	envs [][2]string
}

func (m *Makefile) init(cmd *cobra.Command) {
	cwd := must.NoErrorV(os.Getwd())
	if !IsGoProject(cwd) {
		cmd.Println("please run this command in a go project root")
		os.Exit(1)
	}

	if len(m.Env) > 0 && len(m.Env)%2 != 0 {
		cmd.Println("Env variables must be kv pairs")
		os.Exit(1)
	}

	m.TestIgnore = slicex.Unique(append(DefaultTestIgnores, m.TestIgnore...))
	m.FormatIgnore = slicex.Unique(append(DefaultFormatIgnores, m.FormatIgnore...))

	m.envs = [][2]string{{"GOWORK", "off"}}
	if m.HackTest {
		m.envs = append(m.envs, [2]string{"HACK_TEST", "true"})
	}
	for i := 0; i < len(m.Env); i += 2 {
		m.envs = append(m.envs, [2]string{m.Env[i], m.Env[i+1]})
	}

	m.Depends = append(DefaultDepends, m.Depends...)
	for i := range m.Depends {
		m.Depends[i].SetDefault()
	}
}

func (m *Makefile) Exec(cmd *cobra.Command, args ...string) (err error) {
	m.init(cmd)

	var (
		f *os.File
	)

	f, err = os.OpenFile("Makefile", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer defers.Collect(f.Close, &err)

	m.vars(f)
	m.show(f)
	m.dep(f)
	m.tidy(f)
	m.test(f)
	m.check(f)

	return nil
}

func (m *Makefile) vars(w *os.File) {
	_, _ = fmt.Fprintf(w, `
# go package info
MODULE_PATH    := $(shell cat go.mod | grep ^module -m 1 | awk '{ print $$2; }' || '')
MODULE_NAME    := $(shell basename $(MODULE_PATH))
TEST_IGNORES   := %q
FORMAT_IGNORES := %q

# git repository info
IS_GIT_REPO := $(shell git rev-parse --is-inside-work-tree >/dev/null 2>&1 && echo 1 || echo 0)
ifeq ($(IS_GIT_REPO),1)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
GIT_TAG    := $(shell git describe --tags --abbrev=0)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
else
GIT_COMMIT := ""
GIT_TAG    := ""
GIT_BRANCH := ""
endif
BUILD_AT=$(shell date "+%s")
`, strings.Join(m.TestIgnore, "|"), strings.Join(m.FormatIgnore, ","), "%Y%m%d%H%M%S")

	_, _ = w.WriteString("\n# global env variables\n")
	_ = WriteKeyValAlign(w, "export ", ":=", m.envs)

	_, _ = w.WriteString("\n# go build tools\n")
	_ = WriteKeyValAlign(w, "", ":=", DefaultGoTools)

	kvs := make([][2]string, 0)
	for _, d := range m.Depends {
		kvs = append(kvs, [2]string{
			"DEP_" + stringsx.UpperSnakeCase(d.Name),
			"$(shell type " + d.Name + " > /dev/null 2>&1 && echo $$?)",
		})
	}
	_, _ = w.Write([]byte("\n# dependencies flags\n"))
	_ = WriteKeyValAlign(w, "", ":=", kvs)
}

func (m *Makefile) show(w *os.File) {
	text := `
show:
	@echo "module:"
	@echo "	path=$(MODULE_PATH)"
	@echo "	module=$(MODULE_NAME)"
	@echo "git:"
	@echo "	commit_id=$(GIT_COMMIT)"
	@echo "	tag=$(GIT_TAG)"
	@echo "	branch=$(GIT_BRANCH)"
	@echo "	build_time=$(BUILD_AT)"
	@echo "	name=$(MODULE_NAME)"
	@echo "tools:"
	@echo "	build=$(GOBUILD)"
	@echo "	test=$(GOTEST)"
`
	for _, d := range m.Depends {
		text += fmt.Sprintf(`	@echo "	%s=$(shell which %s)"`, d.Name, d.Name) + "\n"
		text += fmt.Sprintf(`	@echo "	%s=$(%s)"`, d.DepKey(), d.DepKey()) + "\n"
	}
	_, _ = w.WriteString(text)
}

func (m *Makefile) dep(w *os.File) {
	text := `
dep:
	@echo "==> installing dependencies"`
	for _, d := range m.Depends {
		text += fmt.Sprintf(`
	@if [ "${%s}" != "0" ]; then \
		echo "	%s for %s"; \
		go install %s@%s; \
	fi`, d.DepKey(), d.Name, d.Description, d.Repo, d.Version)
	}

	_, _ = w.WriteString(text + "\n")

	text = `
upgrade-dep:
	@echo "==> upgrading dependencies"`

	for _, d := range m.Depends {
		text += fmt.Sprintf(`
	@echo "	%s for %s"
	@go install %s@%s`, d.Name, d.Description, d.Repo, d.Version)
	}

	_, _ = w.WriteString(text + "\n")
}

func (m *Makefile) tidy(w *os.File) {
	text := `
tidy:
	@echo "==> installing dependencies"
	@go mod tidy`
	_, _ = w.WriteString(text + "\n")
}

func (m *Makefile) test(w *os.File) {
	m.hack(w)

	dep := "dep tidy"
	if m.HackTest {
		dep += " hack_dep_run"
	}

	text := fmt.Sprintf(`
test: %s
	@echo "==> run unit test"
	@$(GOTEST) test ./... -race -failfast -parallel 1 -gcflags="all=-N -l"

cover: %s
	@echo "==> run unit test with coverage"
	@$(GOTEST) test ./... -failfast -parallel 1 -gcflags="all=-N -l" -covermode=count -coverprofile=cover.out
	@grep -vE $(TEST_IGNORES) cover.out > cover2.out && mv cover2.out cover.out

view-cover: cover
	@echo "==> run unit test with coverage and view results"
	@$(GOBUILD) tool cover -html cover.out

ci-cover: cover
`, dep, dep)
	_, _ = w.WriteString(text + "\n")
}

func (m *Makefile) hack(w *os.File) {
	if !m.HackTest {
		return
	}

	_, _ = w.WriteString(`
hack_dep_run:
	@cd hack && docker compose up -d --remove-orphans

hack_dep_stop:
	@cd hack && docker compose down -v
`)
}

func (m *Makefile) check(w *os.File) {
	text := `
fmt: dep clean
	@echo "==> formating code"
	@goimports-reviser -rm-unused \
		-imports-order 'std,general,company,project' \
		-project-name ${MODULE_PATH} \
		-excludes $(FORMAT_IGNORES) ./...

lint: dep
	@echo "==> linting"
	@echo ">>>golangci-lint"
	@golangci-lint run
	@echo "done"

clean:
	@find . -name cover.out | xargs rm -rf
	@find . -name .xgo | xargs rm -rf
	@rm -rf build/*

chglog:
	@git chglog -o CHANGELOG.md

pre-commit: dep fmt lint view-cover chglog
`
	_, _ = w.WriteString(text)
}
