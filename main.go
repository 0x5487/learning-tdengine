package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	_ "github.com/taosdata/driver-go/v3/taosSql"

	"github.com/nite-coder/blackbear/pkg/log"
	"github.com/nite-coder/blackbear/pkg/log/handler/console"
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
	logger := log.New()
	clog := console.New()
	logger.AddHandler(clog, log.AllLevels...)
	log.SetLogger(logger)

	pwd, _ := os.Getwd()
	metadata := []Metadata{
		{
			Filename: "trades_btc_usdt_gate",
			Market:   "BTC/USDT",
		},
		// {
		// 	Filename: "gate_trades_eth_usdt",
		// 	Market:   "ETH/USDT",
		// },
	}

	for _, metainfo := range metadata {
		path := filepath.Join(pwd, "test", metainfo.Filename+".csv")
		trades := Load(path)
		log.Infof("filename: %s, count: %d ", metainfo.Filename, len(trades))
		_ = insert(metainfo, trades)
	}

}

func insert(metadata Metadata, data []*Trade) error {
	var taosDSN = "root:taosdata@tcp(host.docker.internal:6030)/demo"
	conn, err := sql.Open("taosSql", taosDSN)
	if err != nil {
		log.Err(err).Error("can't conn to tdengine")
		return err
	}
	defer conn.Close()
	log.Info("connected.")

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

		sb.WriteString(fmt.Sprintf("(%d, %s, %f, %f, %f, %d),", trade.Ts.UnixMicro(), trade.OrderID, trade.Price.BigFloat(), trade.Size.BigFloat(), trade.Price.Mul(trade.Size).BigFloat(), trade.Side))

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
		log.Fatalf("找不到CSV檔案路徑:", path, err)
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
			log.Fatalf("%v", err)
		}

		f, err := strconv.ParseFloat(strings.TrimSpace(record[0]), 64)
		if err != nil {
			panic(err)
		}

		t := f * 1000000

		side, err := strconv.Atoi(strings.TrimSpace(record[4]))
		if err != nil {
			panic(err)
		}

		trade := &Trade{
			Ts:      time.UnixMicro(int64(t)),
			OrderID: strings.TrimSpace(record[1]),
			Price:   decimal.RequireFromString(strings.TrimSpace(record[2])),
			Size:    decimal.RequireFromString(strings.TrimSpace(record[3])),
			Side:    int8(side),
		}
		result = append(result, trade)
	}

	return result
}
