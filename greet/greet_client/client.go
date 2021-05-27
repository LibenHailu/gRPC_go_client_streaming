package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/LibenHailu/grpc-go/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello i am a client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not found connect %v ", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	// fmt.Printf("created client : %f", c)
	// doUnary(c)

	// doServerStreaming(c)

	doClientStreaming(c)

}

func doUnary(c greetpb.GreetServiceClient) {

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Liben",
			LastName:  "Hailu",
		},
	}
	res, err := c.Greet(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}

	log.Printf("Response form Great: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a server streaming RPC")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Liben",
			LastName:  "Hailu",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GetManyTime RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()

		if err == io.EOF {
			// we've reached the end of stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v ", msg.GetResult())

	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a client streaming RPC")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Liben",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Hailu",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Taddessa",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Wagga",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Jijo",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())

	if err != nil {
		log.Fatalf("error while calling  LongGreet: %v", err)
	}

	// we iterate over slice  and send each message individually
	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while reciving response from LongGreet: %v", err)

	}

	fmt.Printf("LongGreet Response: %v\n", res)
}
