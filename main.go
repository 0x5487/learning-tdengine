package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/taosdata/driver-go/v2/af"
	_ "github.com/taosdata/driver-go/v2/taosSql"
)

type Trade struct {
	Ts      time.Time
	OrderID string
	Price   decimal.Decimal
	Size    decimal.Decimal
	Side    int8
}

type Metadata struct {
	Filename string
	Market   string
}

func main() {
	pwd, _ := os.Getwd()
	metadata := []Metadata{
		// {
		// 	Filename: "gate_trades_btc_usdt",
		// 	Market:   "BTC/USDT",
		// },
		{
			Filename: "gate_trades_eth_usdt",
			Market:   "ETH/USDT",
		},
	}

	for _, metainfo := range metadata {
		path := filepath.Join(pwd, "test", metainfo.Filename+".csv")
		trades := Load(path)
		fmt.Println("filename: ", metainfo.Filename, " count: ", len(trades))
		_ = insert(metainfo, trades)
	}

}

func insert(metadata Metadata, data []*Trade) error {
	conn, err := af.Open("host.docker.internal", "root", "taosdata", "demo", 6030)
	defer conn.Close()

	if err != nil {
		fmt.Println("failed to connect, err:", err)
	} else {
		fmt.Println("connected")
	}

	// batch insert
	sb := strings.Builder{}
	for idx, trade := range data {

		if idx%10000 == 0 {
			fmt.Println("insert: ", idx)
			sb = strings.Builder{}
			stmt := fmt.Sprintf("insert into %s USING trades TAGS ('%s', 'gate') values ", metadata.Filename, metadata.Market)
			sb.WriteString(stmt)
		}

		// insert into gate_trades_eth_usdt USING trades TAGS ('ETH/USDT', 'gate') values  (1654041597002, 3620260870, 31793.520000, 0.000800, 1),(1654041589001, 3620260049, 31782.200000, 0.289800, 2);
		//sb.WriteString(fmt.Sprintf("('%s', %s, %f, %f, %d),", trade.Ts.Format("2006-01-02 15:04:05"), trade.OrderID, trade.Price.BigFloat(), trade.Size.BigFloat(), trade.Side))
		//unixMilli := trade.Ts.UnixMilli() + int64(mi)

		sb.WriteString(fmt.Sprintf("(%d, %s, %f, %f, %d),", trade.Ts.UnixMilli(), trade.OrderID, trade.Price.BigFloat(), trade.Size.BigFloat(), trade.Side))

		if idx%10000 == 9999 || idx == len(data)-1 {
			sql := sb.String()
			sql = sql[:len(sql)-1]
			sql += ";"
			_, err := conn.Exec(sql)
			if err != nil {
				fmt.Println("failed to insert, err:", err)
			}
		}
	}

	return nil
}

func Load(path string) []*Trade {
	file, err := os.OpenFile(path, os.O_RDONLY, 0777) // os.O_RDONLY 表示只讀、0777 表示(owner/group/other)權限
	if err != nil {
		log.Fatalln("找不到CSV檔案路徑:", path, err)
	}

	result := make([]*Trade, 0, 1000000)

	// read
	r := csv.NewReader(file)
	r.Comma = ',' // 以何種字元作分隔，預設為`,`。所以這裡可拿掉這行
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		t := strings.Split(strings.TrimSpace(record[0]), ".")

		t0, err := strconv.ParseInt(t[0], 10, 64)
		if err != nil {
			panic(err)
		}
		t1, err := strconv.ParseInt(t[1], 10, 64)
		if err != nil {
			panic(err)
		}
		side, err := strconv.Atoi(strings.TrimSpace(record[4]))
		if err != nil {
			panic(err)
		}

		trade := &Trade{
			Ts:      time.Unix(t0, t1),
			OrderID: record[1],
			Price:   decimal.RequireFromString(strings.TrimSpace(record[2])),
			Size:    decimal.RequireFromString(strings.TrimSpace(record[3])),
			Side:    int8(side),
		}
		result = append(result, trade)
	}

	return result
}
