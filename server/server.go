package main

import (
	"context"

	"log"
	"net"
	"os"
	"strconv"

	pb "github.com/MalteBlackN/30dec/proto"

	"google.golang.org/grpc"
)

type Server struct {
	pb.AuctionServiceServer
}

var currentHighestBid int32 = 0
var HighestBidder string = ""

func (*Server) Bid(ctx context.Context, req *pb.BidRequest) (*pb.BidResponse, error) {
	log.Printf("Received bid request from %v with bid %v", req.GetName(), req.GetBid())
	if req.GetBid() > currentHighestBid {
		currentHighestBid = req.GetBid()
		HighestBidder = req.GetName()
		log.Printf("Bid accepted, new highest bid is %v", currentHighestBid)
		return &pb.BidResponse{Success: true}, nil
	} else {
		log.Printf("Bid rejected, highest bid is %v", currentHighestBid)
		return &pb.BidResponse{Success: false, HighestBid: currentHighestBid, HighestBidder: HighestBidder}, nil
	}
}

func (*Server) Result(ctx context.Context, req *pb.ResultRequest) (*pb.ResultResponse, error) {
	log.Printf("Received result request")
	return &pb.ResultResponse{HighestBid: currentHighestBid}, nil
}

func main() {

	//unlock
	//lock <- true

	//Setting portnumber
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	ownPort := int32(arg1) + 9080
	ownPortStr := strconv.Itoa(int(ownPort))
	log.Println("Starting server on port " + ownPortStr)

	//Listening on own port and creating and setting up server
	list, err := net.Listen("tcp", ":"+ownPortStr)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", ownPortStr, err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterAuctionServiceServer(grpcServer, &Server{})
	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
