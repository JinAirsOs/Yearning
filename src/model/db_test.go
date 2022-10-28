package model

import (
	"gorm.io/gorm"
	"sync"
	"testing"
)

func TestJsonGet(t *testing.T) {
	//var a GlobalConfiguration
	//var c Inception
	//JsonGet(&a)
	//if err := json.Unmarshal(a.Inception, &c); err != nil {
	//
	//}
	//fmt.Println(c.BackUser)
}

func Mock(db *gorm.DB, wg *sync.WaitGroup) {
	for i := 0; i < 100000; i++ {
		db.Create(&CoreDataSource{
			IDC:      "321",
			Username: "3211",
			Password: "cxnudchsj",
			Port:     2,
			Source:   "cnnskj",
			IP:       "156237123032810",
			IsQuery:  1,
		})
	}
	wg.Done()
}
