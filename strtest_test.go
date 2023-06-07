package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	MaxWorkers   = 10                                                               // 并行插入的最大协程数
	TotalRows    = 1000000                                                          // 要插入的总行数
	RandomString = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890" // 用于生成随机字符串的字符集

)

var cities = []string{
	"北京", "上海", "广州", "深圳", "杭州",
	"南京", "成都", "重庆", "武汉", "西安",
	"苏州", "青岛", "天津", "厦门", "长沙",
	"大连", "宁波", "郑州", "沈阳", "济南",
	"哈尔滨", "福州", "合肥", "无锡", "昆明",
	"南宁", "长春", "温州", "佛山", "南昌",
}

func TestVictor(t *testing.T) {
	db, err := sql.Open("mysql", "utest:paswd123@tcp(10.128.23.41:3306)/dbtest")
	if err != nil {
		log.Fatal(err)

	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)

	}

	startTime := time.Now()

	var wg sync.WaitGroup
	jobs := make(chan int, MaxWorkers)

	// 创建并行插入协程
	for i := 0; i < MaxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				row := <-jobs
				if row == -1 { // -1表示没有更多的任务
					return
				}
				err := insertData(db, row)
				if err != nil {
					log.Println(err)

				}

			}

		}()

	}

	// 分配任务给并行插入协程
	for i := 1; i <= TotalRows; i++ {
		jobs <- i

	}
	close(jobs)

	wg.Wait()

	elapsed := time.Since(startTime)
	fmt.Printf("插入 %d 行数据完成，耗时: %s\n", TotalRows, elapsed)

}

// 插入数据到数据库
func insertData(db *sql.DB, row int) error {
	c1 := row
	c2 := generateRandomString(30)
	c3 := time.Now().Format("2006-01-02 15:04:05.000")
	c4 := getRandomCity()

	_, err := db.Exec("INSERT INTO tbstress (c1, c2, c3, c4) VALUES (?, ?, ?, ?)", c1, c2, c3, c4)
	return err

}

// 生成指定长度的随机字符串
func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = RandomString[rand.Intn(len(RandomString))]

	}
	return string(b)

}

// 获取随机中国城市名称
func getRandomCity() string {
	return cities[rand.Intn(len(cities))]

}
