package app

import (
	"context"
	"fmt"
	serverPb "github.com/DaniilOr/gRPC2/pkg/server"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)
type Server struct {
	db1  *pgxpool.Pool
	db2 *pgxpool.Pool
	db3 *pgxpool.Pool
	ctx context.Context
}

func NewServer(db1 *pgxpool.Pool, db2 *pgxpool.Pool, db3 *pgxpool.Pool, ctx context.Context) *Server {
	return &Server{db1: db1, db2: db2, db3: db3, ctx: ctx}
}

func getFromDB(db *pgxpool.Pool, ctx context.Context, req *serverPb.GetRequest)(*serverPb.GetResponse, error){
	var response serverPb.Flight
	var responses []*serverPb.Flight
	date := req.Date.AsTime()
	rows, err := db.Query(ctx, `
	SELECT * FROM flights WHERE fromCity=$1 AND toCity=$2 AND departureTime=$3 LIMIT 50 `,
	req.From, req.To, fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day()))
	if err != nil{
		log.Println(err)
		return nil, err
	}
	for rows.Next() {
		var flightId int64
		var from string
		var to string
		var date time.Time
		rows.Scan(
			&flightId,
			&from,
			&to,
			&response.Duration,
			&response.Cost,
			&date,
		)
		response.DepTime = timestamppb.New(date)
		responses = append(responses, &response)
	}
	if rows.Err() != nil{
		log.Println(rows.Err())
		return nil, rows.Err()
	}
	return &serverPb.GetResponse{Items: responses}, nil
}

func (s*Server) Get(request *serverPb.GetRequest, server serverPb.AgrigatorService_GetServer) error{
	c := make(chan *serverPb.GetResponse)
	errc := make(chan error)
	//done := make(chan error)
	go func(){
		resp, err := getFromDB(s.db1, s.ctx, request)
		if err != nil{
			errc <- err
			return
		}
		errc <- nil
		c<-resp
	}()
	go func(){
		resp, err := getFromDB(s.db2, s.ctx, request)
		if err != nil{
			errc <- err
			return
		}
		errc <- nil
		c<-resp
	}()
	go func(){
		resp, err := getFromDB(s.db3, s.ctx, request)
		if err != nil{
			errc <- err
			return
		}
		errc <- nil
		c<-resp
	}()
	for i := 0; i < 3; i++{
		err := <-errc
		if err == nil {
			response := <-c
			if err := server.Send(response); err != nil {
				log.Print(err)
				if server.Context().Err() == context.DeadlineExceeded {
					log.Print("ctx canceled")
				}
				return err
			}
		} else {
			log.Println(err)
			return err
		}
	}
	return nil
}
