package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"time"

	"github.com/sinmetal/til/grpc/quickstart/echo/echo"
	"golang.org/x/net/context"
	"google.golang.org/api/idtoken"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	grpcMetadata "google.golang.org/grpc/metadata"
	//"google.golang.org/grpc/codes"
	//healthpb "google.golang.org/grpc/health/grpc_health_v1"
	//"google.golang.org/grpc/status"
)

const (
	audience = "244578309736-7m6gpcggkso1ut5l6bsc1lgdot5p0uq9.apps.googleusercontent.com"
)

var (
	conn *grpc.ClientConn
)

func main() {
	ctx := context.Background()

	address := flag.String("host", "localhost:50051", "host:port of gRPC server")
	insecure := flag.Bool("insecure", false, "connect without TLS")
	flag.Parse()

	var err error
	if *insecure == true {
		conn, err = grpc.Dial(*address, grpc.WithInsecure())
	} else {

		systemRoots, err := x509.SystemCertPool()
		if err != nil {
			panic(err)
		}
		cred := credentials.NewTLS(&tls.Config{
			RootCAs: systemRoots,
		})
		conn, err = grpc.Dial(*address, grpc.WithTransportCredentials(cred))
	}

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := echo.NewEchoServerClient(conn)

	/*
		ctx, cancel := context.WithTimeout(ctx, 1 * time.Second)
		defer cancel()
		resp, err := healthpb.NewHealthClient(conn).Check(ctx, &healthpb.HealthCheckRequest{Service: "echo.EchoServer"})
		if err != nil {
			log.Fatalf("HealthCheck failed %+v", err)
		}

		if resp.GetStatus() != healthpb.HealthCheckResponse_SERVING {
			log.Fatalf("service not in serving state: ", resp.GetStatus().String())
		}
		log.Printf("RPC HealthChekStatus:%v", resp.GetStatus())
	*/

	// Create an identity token.
	// With a global TokenSource tokens would be reused and auto-refreshed at need.
	// A given TokenSource is specific to the audience.
	tokenSource, err := idtoken.NewTokenSource(ctx, audience)
	if err != nil {
		panic(err)
	}
	token, err := tokenSource.Token()
	if err != nil {
		panic(err)
	}

	// Add token to gRPC Request.
	ctx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token.AccessToken)

	for i := 0; i < 10; i++ {
		r, err := c.SayHello(ctx, &echo.EchoRequest{Name: "unary RPC msg "})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		time.Sleep(1 * time.Second)
		log.Printf("RPC Response: %v %v", i, r)
	}

	/*
		stream, err := c.SayHelloStream(ctx, &pb.EchoRequest{Name: "Stream RPC msg"}, grpc.Header(&header))
		if err != nil {
			log.Fatalf("SayHelloStream(_) = _, %v", err)
		}
		for {
			m, err := stream.Recv()
			if err == io.EOF {
				t := stream.Trailer()
				log.Println("Stream Trailer: ", t)
				break
			}
			if err != nil {
				log.Fatalf("SayHelloStream(_) = _, %v", err)
			}
			h, err := stream.Header()
			if err != nil {
				log.Fatalf("stream.Header error _, %v", err)
			}
			log.Printf("Stream Header: ", h)
			log.Printf("Message: ", m.Message)
		}
	*/

}
