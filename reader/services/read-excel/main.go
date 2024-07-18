package main

import (
	"log"
	"net"
	"time"

	read_excel "github.com/sonlam137216/read-excel/services/common/genproto/read-excel"
	"github.com/thedatashed/xlsxreader"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type server struct {
	read_excel.UnimplementedExampleServiceServer
}

func (s *server) SayHello(ctx context.Context, req *read_excel.HelloRequest) (*read_excel.HelloResponse, error) {
	startTime := time.Now()

	xl, err := xlsxreader.OpenFile("./file.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer xl.Close()

	rowData := make([]*read_excel.Row, 0)

	for row := range xl.ReadRows(xl.Sheets[0]) {
		cellData := make([]*read_excel.Cell, 0)
		for _, cell := range row.Cells {
			cellData = append(cellData, &read_excel.Cell{
				Column: cell.Column,
				Row:    int32(cell.Row),
				Value:  cell.Value,
			})
		}
		rowData = append(rowData, &read_excel.Row{Cells: cellData})
	}

	endTime := time.Now()

	// Log the duration
	duration := endTime.Sub(startTime)
	log.Printf("Operation duration: %s\n", duration)

	return &read_excel.HelloResponse{Rows: rowData}, nil
}

func main() {

	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	read_excel.RegisterExampleServiceServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve: %v", err)
	}
}
