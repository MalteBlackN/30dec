package main

import (
	"context"

	"log"

	"os"
	"strconv"
	"strings"

	"bufio"

	pb "github.com/MalteBlackN/30dec/proto"

	"google.golang.org/grpc"
)

var name string

var totalPorts int64
var reader = bufio.NewReader(os.Stdin)

var clients []pb.AuctionServiceClient

func main() {
	//Loading id and total amount of ports to connect to

	totalPorts, _ = strconv.ParseInt(os.Args[1], 10, 32)

	//Creating connection to all servers
	for i := 0; i < int(totalPorts); i++ {
		// Create a virtual RPC Client Connection on port 9080 + i
		var conn *grpc.ClientConn
		var port int = 9080 + i
		portStr := strconv.Itoa(port)

		conn, err := grpc.Dial(":"+portStr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		// Defer means: When this function returns, call this method (meaning, one main is done, close connection)
		defer conn.Close()

		//  Create new Client from generated gRPC code from proto
		c := pb.NewAuctionServiceClient(conn)
		clients = append(clients, c)
	}

	log.Print("wlcome to the the auction, please enter your name, followed by enter:")
	name, _ = reader.ReadString('\n')
	name = strings.TrimSpace(name)

	//Starting method for continuously recieving input from user
	for {
		takeInput()
	}
}

func takeInput() {
	for {

		log.Println("Write 'bid' to bid on an item, 'result' to see the current highest bid, followed by enter:")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		log.Print("(You wrote: " + input + ")")

		if input == "bid" {
			log.Println("write the amount you want to bid, followed by enter:")

			tempAmount, _ := reader.ReadString('\n')
			tempAmount = strings.TrimSpace(tempAmount)
			amount, err := strconv.Atoi(tempAmount)
			//log.Println(("you wish to bid: %d on the item"), amount)
			if err != nil {
				log.Println("Faulty input, please try again")
				continue
			}

			//Calling method to bid the amount on the item
			result, _ := bidAmountForAll(&pb.BidRequest{Name: name, Bid: int32(amount)})

			//Displaying information to user based on whether the value was put in the hash table or not
			if result.Success {
				log.Printf(("your bid of: %d was place, you are now leading the auction)\n"), amount)
				continue
			} else {
				log.Printf(("your bid is not high enough!! %s holds the highest bid of: %d)\n"), result.HighestBidder, result.HighestBid)
				break
			}
		}
		if input == "result" {

			log.Println("the current highest bid is:")

			//Calling method to get value from hash table
			result, _ := getResultFromAll(&pb.ResultRequest{})

			if result.Success {
				log.Printf("the value could not be retrieved)\n")
				continue
			} else {
				log.Printf(("the current highest bid is: %d"), result.HighestBid)
				break
			}

		}
	}
}

func bidAmountForAll(req *pb.BidRequest) (*pb.BidResponse, error) {
	var ack *pb.BidResponse
	for _, c := range clients {
		tempAck, err := c.Bid(context.Background(), req)
		if err == nil {
			ack = tempAck
		}
	}
	return ack, nil
}

func getResultFromAll(req *pb.ResultRequest) (*pb.ResultResponse, error) {
	var result *pb.ResultResponse
	for _, c := range clients {
		tempResult, err := c.Result(context.Background(), req)
		if err == nil {
			result = tempResult
		}
	}
	return result, nil
}
