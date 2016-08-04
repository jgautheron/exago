package repository

import "testing"

func TestGotData(t *testing.T) {
	rp, _ := loadStubRepo()
	data := rp.GetData()
	if data.CodeStats["CLOC"] != 48 {
		t.Error("The data could not be retrieved")
	}
}

func TestGotName(t *testing.T) {
	rp, _ := loadStubRepo()
	if rp.GetName() != repo {
		t.Error("The repository name is wrong")
	}
}

func TestGotRank(t *testing.T) {
	rp, _ := loadStubRepo()
	if rp.GetRank() != "F" {
		t.Error("The rank is wrong")
	}
}
