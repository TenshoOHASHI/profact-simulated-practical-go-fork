package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func main() {
	file, err := os.Create("temp/test_properties_utf8.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// ヘッダー
	writer.Write([]string{"物件名", "賃料", "住所", "間取り", "ステータス"})

	layouts := []string{"1K", "1DK", "1LDK", "2K", "2DK", "2LDK", "3K", "3DK", "3LDK"}
	statuses := []string{"available", "contracted", "hidden"}

	for i := 1; i <= 10000; i++ {
		name := fmt.Sprintf("テスト物件%d", i)
		rent := 50000 + (i % 100000)
		address := fmt.Sprintf("東京都渋谷区%d-%d-%d", i%100, i%50, i%30)
		layout := layouts[i%len(layouts)]
		status := statuses[i%len(statuses)]

		writer.Write([]string{
			name,
			fmt.Sprintf("%d", rent),
			address,
			layout,
			status,
		})
	}
}
