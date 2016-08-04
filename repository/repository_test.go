package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/exago/svc/mocks"
)

var (
	repo = "github.com/foo/bar"
	data = `{"codestats":{"Assertion":0,"CLOC":48,"ExportedFunction":9,"ExportedMethod":6,"Function":15,"FunctionLOC":197,"GoStatement":0,"IfStatement":28,"Import":19,"Interface":2,"LOC":285,"Method":6,"MethodLOC":6,"NCLOC":237,"NOF":3,"Struct":3,"SwitchStatement":1,"Test":1},"imports":[],"testresults":{"checklist":{"Failed":[{"Category":"extraCredit","Desc":"Blackbox Tests: In addition to standard tests, does the project have blackbox tests?","Name":"hasBlackboxTests"},{"Category":"extraCredit","Desc":"Benchmarks: In addition to tests, does the project have benchmarks?","Name":"hasBenches"}],"Passed":[{"Category":"minimumCriteria","Desc":"README Presence: Does the project's include a documentation entrypoint?","Name":"hasReadme"},{"Category":"minimumCriteria","Desc":"Directory Names and Packages Match: Does each package <pkg> statement's package name match the containing directory name?","Name":"isDirMatch"},{"Category":"goodCitizen","Desc":"Contribution Process: Does the project document a contribution process?","Name":"hasContributing"},{"Category":"minimumCriteria","Desc":"Licensed: Does the project have a license?","Name":"hasLicense"},{"Category":"minimumCriteria","Desc":"golint Correctness: Is the linter satisfied?","Name":"isLinted"},{"Category":"minimumCriteria","Desc":"gofmt Correctness: Is the code formatted correctly?","Name":"isFormatted"},{"Category":"minimumCriteria","Desc":"go tool vet Correctness: Is the Go vet satisfied?","Name":"isVetted"}]},"packages":[{"coverage":42.9,"execution_time":0.005,"name":"github.com/jgautheron/codename-generator","success":true,"tests":[{"name":"TestCodenameGeneration","execution_time":0,"passed":true}]}],"execution_time":{"goget":"1.594812149s","goprove":"254.995279ms","gotest":"714.61787ms"},"raw_output":{"goget":"","gotest":"=== RUN   TestCodenameGeneration\n--- PASS: TestCodenameGeneration (0.00s)\nPASS\ncoverage: 42.9% of statements\nok  \tgithub.com/jgautheron/codename-generator\t0.005s\n"},"errors":{"goget":"","gotest":""}},"lintmessages":{"codename.go":{"errcheck":[{"col":16,"line":45,"message":"error return value not checked (json.Unmarshal(sp, &spa))","severity":"warning"},{"col":16,"line":53,"message":"error return value not checked (json.Unmarshal(sh, &sha))","severity":"warning"}],"golint":[{"col":2,"line":29,"message":"exported const SuperbFilePath should have comment (or a comment on this block) or be unexported","severity":"warning"},{"col":6,"line":33,"message":"exported type FormatType should have comment or be unexported","severity":"warning"},{"col":6,"line":34,"message":"exported type JSONData should have comment or be unexported","severity":"warning"},{"col":1,"line":36,"message":"comment on exported function Get should be of the form \"Get ...\"","severity":"warning"}]},"words.go":{"dupl":[{"col":0,"line":81,"message":"duplicate of words.go:101-110","severity":"warning"},{"col":0,"line":101,"message":"duplicate of words.go:81-90","severity":"warning"},{"col":0,"line":115,"message":"duplicate of words.go:141-151","severity":"warning"},{"col":0,"line":141,"message":"duplicate of words.go:115-125","severity":"warning"}],"gofmt":[{"col":0,"line":1,"message":"file is not gofmted","severity":"warning"}],"golint":[{"col":5,"line":72,"message":"var _dataSuperbJson should be _dataSuperbJSON","severity":"warning"},{"col":6,"line":74,"message":"func dataSuperbJsonBytes should be dataSuperbJSONBytes","severity":"warning"},{"col":6,"line":81,"message":"func dataSuperbJson should be dataSuperbJSON","severity":"warning"},{"col":5,"line":92,"message":"var _dataSuperheroesJson should be _dataSuperheroesJSON","severity":"warning"},{"col":6,"line":94,"message":"func dataSuperheroesJsonBytes should be dataSuperheroesJSONBytes","severity":"warning"},{"col":6,"line":101,"message":"func dataSuperheroesJson should be dataSuperheroesJSON","severity":"warning"}],"gosimple":[{"col":2,"line":234,"message":"'if err != nil { return err }; return nil' can be simplified to 'return err'","severity":"warning"}]}},"metadata":{"image":"https://avatars.githubusercontent.com/u/683888?v=3","description":"A codename generator meant for naming software releases.","stars":13,"last_push":"2015-08-29T20:32:12Z"},"score":{"value":68.72365160829756,"rank":"F","details":[{"name":"imports","score":100,"weight":1.5,"desc":"counts the number of third party libraries","msg":"0 third-party package(s)","url":"https://github.com/jgautheron/gogetimports"},{"name":"testcoverage","score":37.2941129735386,"weight":2.45484486000851,"desc":"measures pourcentage of code covered by tests","msg":"coverage is greater or equal to 42.90","url":"https://golang.org/pkg/testing/"},{"name":"checklist","score":94.56521739130436,"weight":1.8,"desc":"inspects project for the best practices listed in the Go CheckList","msg":"","url":"https://github.com/karolgorecki/goprove","details":[{"name":"hasBenches","score":0,"weight":0.5,"msg":"check failed"},{"name":"projectBuilds","score":100,"weight":1.5,"msg":"check succeeded","url":"https://github.com/matttproud/gochecklist/blob/master/publication/compilation.md"},{"name":"isFormatted","score":100,"weight":3,"msg":"check succeeded","url":"https://github.com/matttproud/gochecklist/blob/master/publication/code_correctness.md"},{"name":"hasReadme","score":100,"weight":3,"msg":"check succeeded","url":"https://github.com/matttproud/gochecklist/blob/master/publication/documentation_entrypoint.md"},{"name":"isDirMatch","score":100,"weight":0.7,"msg":"check succeeded","url":"https://github.com/matttproud/gochecklist/blob/master/publication/dir_pkg_match.md"},{"name":"isLinted","score":100,"weight":0.5,"msg":"check succeeded","url":"https://github.com/matttproud/gochecklist/blob/master/publication/code_correctness.md"}]},{"name":"codestats","score":8.274316016134177,"weight":1,"desc":"counts lines of code, comments, functions, structs, imports etc in Go code","msg":"48 comments for 285 lines of code","url":"https://github.com/jgautheron/golocc"},{"name":"testduration","score":99.98571549999274,"weight":1.2,"desc":"measures time taken for testing","msg":"tests took 0.01s","url":"https://golang.org/pkg/testing/"},{"name":"lintmessages","score":72.05373125586767,"weight":2,"desc":"runs a whole bunch of Go linters","msg":"","url":"https://github.com/alecthomas/gometalinter","details":[{"name":"gofmt","score":53.17515301305708,"weight":3,"desc":"detects if Go code is incorrectly formatted","msg":"exceeds the warnings/LOC threshold","url":"https://golang.org/cmd/gofmt/"},{"name":"errcheck","score":86.9053250930649,"weight":2,"desc":"finds unchecked errors in Go code","msg":"exceeds the warnings/LOC threshold","url":"https://github.com/kisielk/errcheck"},{"name":"gosimple","score":90.00876262522593,"weight":1.5,"desc":"examines Go code and reports constructs that can be simplified","msg":"exceeds the warnings/LOC threshold","url":"https://github.com/dominikh/go-simple"}]}]},"execution_time":"15s","last_update":"2016-07-12T00:33:25.471167217+02:00"}`
)

func TestIsNotLoaded(t *testing.T) {
	rp := &Repository{
		Name: repo,
	}
	loaded := rp.IsLoaded()
	if loaded {
		t.Errorf("The repository %s should not be loaded", rp.Name)
	}
}

func TestIsLoaded(t *testing.T) {
	rp, err := loadStubRepo()
	if err != nil {
		t.Error(err)
	}

	loaded := rp.IsLoaded()
	if !loaded {
		t.Errorf("The repository %s should be loaded", rp.Name)
	}
}

func TestStartTimeSet(t *testing.T) {
	rp := &Repository{
		Name: repo,
	}
	now := time.Now()
	rp.SetStartTime(now)
	if now != rp.startTime {
		t.Error("Got the wrong time")
	}
}

func loadStubRepo() (*Repository, error) {
	rhMock := mocks.RepositoryHost{}
	dbMock := mocks.Database{}
	dbMock.On("Get",
		[]byte(fmt.Sprintf("%s-%s", repo, "")),
	).Return([]byte(data), nil)

	rp := &Repository{
		Name:           repo,
		DB:             dbMock,
		RepositoryHost: rhMock,
	}
	if err := rp.Load(); err != nil {
		return nil, fmt.Errorf("Got error while loading data: %v", err)
	}
	return rp, nil
}
