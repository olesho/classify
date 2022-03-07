package sequence

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func BenchmarkCrash(t *testing.B) {
	entries, err := ioutil.ReadDir("./rozetka")
	if err != nil {
		t.Error(err)
		return
	}

	cluster := NewRootCluster()
	//cluster.SetLimit(10)
	for _, entry := range entries[:10] {
		if !entry.IsDir() {
			f, err := os.Open("./rozetka/" + entry.Name())
			if err != nil {
				t.Error(err)
			}
			data, _ := ioutil.ReadAll(f)

			err = cluster.LoadString(string(data))
			if err != nil {
				log.Println(err)
				return
			}

			err = f.Close()
			if err != nil {
				t.Error(err)
				return
			}
		}
	}
	cluster.Batch().Results()
}