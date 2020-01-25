package cov

import (
	"os"
	"testing"

	"golang.org/x/tools/cover"
)

var profile = &cover.Profile{
	FileName: "github.com/chreble/todo/task/task.go",
	Mode:     "count",
	Blocks: []cover.ProfileBlock{
		cover.ProfileBlock{18, 30, 20, 25, 2, 0},
		cover.ProfileBlock{23, 2, 23, 14, 1, 0},
		cover.ProfileBlock{20, 25, 22, 3, 1, 0},
		cover.ProfileBlock{27, 38, 28, 33, 1, 0},
		cover.ProfileBlock{31, 2, 31, 12, 1, 0},
		cover.ProfileBlock{28, 33, 30, 3, 1, 0},
		cover.ProfileBlock{35, 40, 39, 2, 3, 1},
		cover.ProfileBlock{42, 58, 43, 37, 1, 0},
		cover.ProfileBlock{46, 2, 46, 12, 1, 0},
		cover.ProfileBlock{43, 37, 45, 3, 1, 0},
		cover.ProfileBlock{50, 35, 51, 19, 1, 0},
		cover.ProfileBlock{54, 2, 54, 11, 1, 0},
		cover.ProfileBlock{51, 19, 53, 3, 1, 0},
		cover.ProfileBlock{58, 41, 59, 19, 1, 0},
		cover.ProfileBlock{64, 2, 64, 11, 1, 0},
		cover.ProfileBlock{59, 19, 60, 14, 1, 0},
		cover.ProfileBlock{60, 14, 62, 4, 1, 0},
		cover.ProfileBlock{68, 21, 71, 2, 2, 1},
		cover.ProfileBlock{74, 40, 78, 2, 3, 0},
	},
}

var mock = struct {
	Coverage float64
	Pkg      Package
}{
	Coverage: 20,
	Pkg: Package{
		Name:     "task",
		Path:     "github.com/chreble/todo/task",
		Coverage: 20,
		Functions: []*Function{
			&Function{"Tasks.All", "$GOPATH/github.com/chreble/todo/task/task.go", 18, 24, 0, 0, nil},
			&Function{"Tasks.Create", "$GOPATH/github.com/chreble/todo/task/task.go", 35, 39, 100, 3, nil},
			&Function{"newID", "$GOPATH/github.com/chreble/todo/task/task.go", 68, 71, 100, 2, nil},
		},
	},
}

var report *Report

func TestMain(m *testing.M) {
	report = &Report{}
	report.parseProfile([]*cover.Profile{profile})

	os.Exit(m.Run())
}

func TestGlobal(t *testing.T) {
	// Check global coverage
	if report.Coverage != mock.Coverage {
		t.Errorf(
			"Got wrong global coverage expected %.2f computed %.2f",
			mock.Coverage,
			report.Coverage,
		)
	}
}

func TestPackage(t *testing.T) {
	// Check Package with a predicate function
	pkg := From(report.Packages).SingleWith(func(p interface{}) bool {
		return p.(*Package).Name == mock.Pkg.Name
	})
	if pkg == nil {
		t.Fatalf("Package %s not found in report", mock.Pkg.Name)
	}
	if pkg.(*Package).Name != mock.Pkg.Name ||
		pkg.(*Package).Coverage != mock.Pkg.Coverage ||
		pkg.(*Package).Path != mock.Pkg.Path {
		// Return information about error
		t.Errorf(
			`Got error on package assertion:
				* Got %s name and expected %s
				* Got %s path and expected %s
				* Got %.2f coverage and expected %.2f
			`,
			pkg.(*Package).Name, mock.Pkg.Name,
			pkg.(*Package).Coverage, mock.Pkg.Coverage,
			pkg.(*Package).Path, mock.Pkg.Path,
		)
	}
}

func TestFunctions(t *testing.T) {
	// Check Functions
	pkg := From(report.Packages).SingleWith(func(p interface{}) bool {
		return p.(*Package).Name == mock.Pkg.Name
	})
	if pkg == nil {
		t.Fatalf("Package %s not found in report", mock.Pkg.Name)
	}
	for _, fn := range pkg.(*Package).Functions {
		// Find function
		mfn := From(mock.Pkg.Functions).SingleWith(func(f interface{}) bool {
			return f.(*Function).Name == fn.Name
		})
		// Ignore un-mocked functions
		if mfn == nil {
			continue
		}
		f := mfn.(*Function)
		if f.Coverage != fn.Coverage {
			t.Errorf(
				"Got %.2f coverage and expected %.2f for Function %s",
				fn.Coverage,
				f.Coverage,
				fn.Name,
			)
		}
	}
}

func TestAccumulatePackage(t *testing.T) {
	p1_1 := registerPackage("p1")
	p1_2 := registerPackage("p1")
	p2 := registerPackage("p2")
	p3 := registerPackage("p1")
	registerFunction(p3, "f", "file.go", 0, 1)
	p4 := registerPackage("p1")
	registerFunction(p4, "f", "file.go", 1, 2)

	var tests = [...]struct {
		a, b       *Package
		expectPass bool
	}{
		// Should work: everything is the same.
		{p1_1, p1_2, true},
		// Should fail: name is different.
		{p1_1, p2, false},
		// Should fail: numbers of functions are different.
		{p1_1, p3, false},
		// Should fail: functions are different.
		{p3, p4, false},
	}

	for _, test := range tests {
		err := test.a.Accumulate(test.b)
		if test.expectPass {
			if err != nil {
				t.Error(err)
			}
		} else {
			if err == nil {
				t.Error("Expected an error")
			}
		}
	}
}

func TestAccumulateFunction(t *testing.T) {
	p := registerPackage("p1")
	f1_1 := registerFunction(p, "f1", "file.go", 0, 1)
	f1_2 := registerFunction(p, "f1", "file.go", 0, 1)
	f2 := registerFunction(p, "f2", "file.go", 0, 1)
	f3 := registerFunction(p, "f1", "file2.go", 0, 1)
	f4 := registerFunction(p, "f1", "file.go", 2, 3)
	f5 := registerFunction(p, "f1", "file.go", 0, 1)
	registerStatement(f5, 0, 1)
	f6 := registerFunction(p, "f1", "file.go", 0, 1)
	registerStatement(f6, 2, 3)

	var tests = [...]struct {
		a, b       *Function
		expectPass bool
	}{
		// Should work: everything is the same.
		{f1_1, f1_2, true},
		// Should fail: names are different.
		{f1_1, f2, false},
		// Should fail: files are different.
		{f1_1, f3, false},
		// Should fail: ranges are different.
		{f1_1, f4, false},
		// Should fail: numbers of statements are different.
		{f1_1, f5, false},
		// Should fail: all the same, except statement values.
		{f5, f6, false},
	}

	for _, test := range tests {
		err := test.a.Accumulate(test.b)
		if test.expectPass {
			if err != nil {
				t.Error(err)
			}
		} else {
			if err == nil {
				t.Error("Expected an error")
			}
		}
	}
}

func TestAccumulateStatement(t *testing.T) {
	p := registerPackage("p1")
	f := registerFunction(p, "f1", "file.go", 0, 1)
	s1_1 := registerStatement(f, 0, 1)
	s1_2 := registerStatement(f, 0, 1)
	s2 := registerStatement(f, 2, 3)

	// Should work: ranges are the same.
	if err := s1_1.Accumulate(s1_2); err != nil {
		t.Error(err)
	}

	// Should fail: ranges are not the same.
	if err := s1_1.Accumulate(s2); err == nil {
		t.Errorf("Expected an error")
	}
}

func registerPackage(name string) *Package {
	return &Package{Name: name}
}

func registerFunction(p *Package, name, file string, startOffset, endOffset int) *Function {
	f := &Function{Name: name, File: file, Start: startOffset, End: endOffset}
	p.Functions = append(p.Functions, f)
	return f
}

func registerStatement(f *Function, startOffset, endOffset int) *Statement {
	s := &Statement{Start: startOffset, End: endOffset}
	f.Statements = append(f.Statements, s)
	return s
}
