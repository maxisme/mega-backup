package mega

import "testing"

func TestBackupDirectory(t *testing.T) {
	mcs := CreateServer{
		Host:        "http://idmy.team:8980",
		Credentials: "GUIGQ31rtwbJIGYS5Jt3syhDxBhYH8uije5WEnnVr2vcBWCfZBAwLRJPLraDDAUfEtQ6gBY5TMkdH6Cl",
	}
	_, err := mcs.BackupDirectory("./test/", "foo.txt", "hello", Account{
		Email:    "1588862414969802312@dlme.ga",
		Password: "orApUc124HpxtStscKS78t26y2ora5Pub7MB6PxghMmLcVie7I",
	})
	if err != nil {
		t.Errorf(err.Error())
	}
}
