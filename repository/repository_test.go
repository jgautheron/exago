package repository

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/exago/svc/mocks"
	"github.com/exago/svc/repository/model"
)

var repo = "github.com/foo/bar"

func TestIsNotCached(t *testing.T) {
	dbMock := mocks.Database{}
	dbMock.On("FindAllForRepository", []byte(
		fmt.Sprintf("%s-%s", repo, "")),
	).Return(map[string][]byte{}, nil)

	rp := &Repository{
		Name: repo,
		db:   dbMock,
	}
	cached := rp.IsCached()
	if cached {
		t.Errorf("The repository %s should not be cached", rp.Name)
	}
}

func TestIsCached(t *testing.T) {
	dbMock := mocks.Database{}
	dbMock.On("FindAllForRepository", []byte(
		fmt.Sprintf("%s-%s", repo, "")),
	).Return(map[string][]byte{
		"codestats":      []byte(""),
		"imports":        []byte(""),
		"testresults":    []byte(""),
		"lintmessages":   []byte(""),
		"metadata":       []byte(""),
		"score":          []byte(""),
		"execution_time": []byte(""),
		"last_update":    []byte(""),
	}, nil)

	rp := &Repository{
		Name: repo,
		db:   dbMock,
	}
	cached := rp.IsCached()
	if !cached {
		t.Errorf("The repository %s should be cached", rp.Name)
	}
}

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
	tr := model.TestResults{}
	tr.RawOutput.Gotest = "foo"

	rp := &Repository{
		Name:         repo,
		CodeStats:    model.CodeStats{"LOC": 10},
		Imports:      []string{"foo"},
		TestResults:  tr,
		LintMessages: model.LintMessages{},
	}
	loaded := rp.IsLoaded()
	if !loaded {
		t.Errorf("The repository %s should be loaded", rp.Name)
	}
}

func TestLoadedFromDB(t *testing.T) {
	rp := getRepositoryMock("A")
	if err := rp.Load(); err != nil {
		t.Errorf("Got error while loading the repo: %v", err)
	}

	if rp.GetName() != repo {
		t.Errorf("The name should be %s", repo)
	}

	desc := "foobar"
	if rp.GetMetadataDescription() != desc {
		t.Errorf("The description should be %s", desc)
	}

	image := "http://foo.com/img.png"
	if rp.GetMetadataImage() != image {
		t.Errorf("The image should be %s", image)
	}

	stars := 99
	if rp.GetMetadataStars() != stars {
		t.Errorf("There should be %d stars", stars)
	}

	if len(rp.Imports) != 2 {
		t.Error("There should be two third parties loaded")
	}

	if rp.GetRank() != "A" {
		t.Error("The rank should be A")
	}
}

func TestCacheCleared(t *testing.T) {
	dbMock := mocks.Database{}
	dbMock.On("DeleteAllMatchingPrefix", []byte(
		fmt.Sprintf("%s-%s", repo, ""),
	)).Return(nil)

	rp := &Repository{
		db:   dbMock,
		Name: repo,
	}
	if err := rp.ClearCache(); err != nil {
		t.Error("Got error while attempting to clear cache")
	}
}

func TestGotMap(t *testing.T) {
	rp := getRepositoryMock("A")
	rp.Load()

	mp := rp.AsMap()
	if mp[model.ScoreName].(model.Score).Rank != "A" {
		t.Errorf("The rank should be A, got %s", mp["Rank"])
	}
}

func TestStartTimeSet(t *testing.T) {
	rp := &Repository{
		Name: repo,
	}
	now := time.Now()
	rp.SetStartTime(now)
	if now != rp.StartTime {
		t.Error("Got the wrong time")
	}
}

func TestScoreCalculated(t *testing.T) {
	rp := getRepositoryMock("")
	rp.Load()
	rp.calcScore()
	if rp.GetRank() == "" {
		t.Error("The rank has not been set")
	}
}

func getRepositoryMock(rank string) *Repository {
	var b []byte
	dbMock := mocks.Database{}

	imports := []string{"foo", "bar"}
	b, _ = json.Marshal(imports)
	dbMock.On("FindForRepositoryCmd", []byte(
		fmt.Sprintf("%s-%s-%s", repo, "", model.ImportsName),
	)).Return(b, nil)

	codestats := map[string]int{"LOC": 123, "NCLOC": 20}
	b, _ = json.Marshal(codestats)
	dbMock.On("FindForRepositoryCmd", []byte(
		fmt.Sprintf("%s-%s-%s", repo, "", model.CodeStatsName),
	)).Return(b, nil)

	lintmessages := map[string]map[string][]map[string]interface{}{}
	b, _ = json.Marshal(lintmessages)
	dbMock.On("FindForRepositoryCmd", []byte(
		fmt.Sprintf("%s-%s-%s", repo, "", model.LintMessagesName),
	)).Return(b, nil)

	var testresults model.TestResults
	testresults.RawOutput.Gotest = "foo"
	b, _ = json.Marshal(testresults)
	dbMock.On("FindForRepositoryCmd", []byte(
		fmt.Sprintf("%s-%s-%s", repo, "", model.TestResultsName),
	)).Return(b, nil)

	var score model.Score
	score.Rank = rank
	b, _ = json.Marshal(score)
	dbMock.On("FindForRepositoryCmd", []byte(
		fmt.Sprintf("%s-%s-%s", repo, "", model.ScoreName),
	)).Return(b, nil)

	executionTime := "20s"
	b, _ = json.Marshal(executionTime)
	dbMock.On("FindForRepositoryCmd", []byte(
		fmt.Sprintf("%s-%s-%s", repo, "", model.ExecutionTimeName),
	)).Return(b, nil)

	var metadata model.Metadata
	metadata.Description = "foobar"
	metadata.Image = "http://foo.com/img.png"
	metadata.Stars = 99
	b, _ = json.Marshal(metadata)
	dbMock.On("FindForRepositoryCmd", []byte(
		fmt.Sprintf("%s-%s-%s", repo, "", model.MetadataName),
	)).Return(b, nil)

	lastUpdate := time.Now()
	b, _ = json.Marshal(lastUpdate)
	dbMock.On("FindForRepositoryCmd", []byte(
		fmt.Sprintf("%s-%s-%s", repo, "", model.LastUpdateName),
	)).Return(b, nil)

	return &Repository{
		db:   dbMock,
		Name: repo,
	}
}
