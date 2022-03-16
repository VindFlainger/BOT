package mymath

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

func GuessChanse(rang, samplecount, pickcount int) float64 {
	var chanse float64 = 1
	freenums := samplecount - pickcount
	for i := 1; i <= freenums; i++ {
		chanse *= (float64(rang-pickcount) - float64(i-1)) / float64(i)
	}

	for i := 1; i <= samplecount; i++ {
		chanse /= (float64(rang) - float64(i-1)) / float64(i)
	}

	return chanse
}

func GetRandomVal(count, maxVal int) (rands []int) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < count; i++ {
		rannum := rand.Intn(maxVal) + 1
		if func() bool {
			for _, val := range rands {
				if val == rannum {
					return false
				}
			}
			return true
		}() {
			rands = append(rands, rannum)
		} else {
			i--
		}

	}
	return
}

func CheckNums(winnums []int, nums ...int) bool {
	for _, num := range nums {
		for _, winnum := range winnums {
			if num == winnum {
				return true
			}
		}
	}
	return false
}

func CheckAvail(ID int, db *sql.DB) int {
	var points int
	left := db.QueryRow(fmt.Sprintf("SELECT points FROM uir WHERE vk_id = %d; UPDATE uir SET points = points-1 WHERE points > 0 AND vk_id = %d", ID, ID))
	left.Scan(&points)
	return points
}
