package main

import (
	"context"
	"log"
	"time"

	pb "github.com/IM-Deane/mailing-list/proto"
	"github.com/alexflint/go-arg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// logResponse helper that logs a response or error
func logResponse(res *pb.EmailResponse, err error) {
	if err != nil {
		log.Fatalf(" error: %v", err)
	}

	if res.EmailEntry == nil {
		log.Printf(" email not found")
	} else {
		log.Printf(" response %v", res.EmailEntry)
	}
}


// createEmail handles email creation on client
func createEmail(client pb.MailingListServiceClient, address string) (*pb.EmailEntry) {
	log.Println("create email")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// if request takes less than 1 second we free up resources
	defer cancel()

	res, err := client.CreateEmail(ctx, &pb.CreateEmailRequest{EmailAddr: address})
	logResponse(res, err)

	return res.EmailEntry
}

// getEmail handles fetching an email on client
func getEmail(client pb.MailingListServiceClient, address string) (*pb.EmailEntry) {
	log.Println("get email")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// if request takes less than 1 second we free up resources
	defer cancel()

	res, err := client.GetEmail(ctx, &pb.GetEmailRequest{EmailAddr: address})
	logResponse(res, err)

	return res.EmailEntry
}

// getEmailBatch handles fetching a list of emails on client
func getEmailBatch(client pb.MailingListServiceClient, count int, page int) {
	log.Println("get email batch")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// if request takes less than 1 second we free up resources
	defer cancel()

	res, err := client.GetEmailBatch(ctx, &pb.GetEmailBatchRequest{
		Count: int32(count), Page: int32(page)})
	if err != nil {
		log.Fatalf(" error: %v", err)
	}
	log.Println("response")
	for i := 0; i < len(res.EmailEntry); i++ {
		// print number, total entries, and response
		log.Printf(" item [%v of %v]: %s", i+1, len(res.EmailEntry), res.EmailEntry[i])
	}
}

// updateEmail handles updating emails via client
func updateEmail(client pb.MailingListServiceClient, entry pb.EmailEntry) (*pb.EmailEntry) {
	log.Println("update email")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// if request takes less than 1 second we free up resources
	defer cancel()

	res, err := client.UpdateEmail(ctx, &pb.UpdateEmailRequest{EmailEntry: &entry})
	logResponse(res, err)

	return res.EmailEntry
}

// deleteEmail handles updating emails via client
func deleteEmail(client pb.MailingListServiceClient, address string) (*pb.EmailEntry) {
	log.Println("delete email")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// if request takes less than 1 second we free up resources
	defer cancel()

	res, err := client.DeleteEmail(ctx, &pb.DeleteEmailRequest{EmailAddr: address})
	logResponse(res, err)

	return res.EmailEntry
}

// command line args
var args struct {
	GRPCAddr string `arg:"env:MAILINGLIST_GRPC_ADDR"`
}

func main() {
	arg.MustParse(&args)

	// set default address
	if args.GRPCAddr == "" {
		args.GRPCAddr = ":8081"
	}

	// connect to gRPC server (no encryption as its an internal service)
	conn, err := grpc.Dial(args.GRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewMailingListServiceClient(conn)

	// DEFAULT: get email from DB
	getEmail(client, "example@example.com")

	// *** UNCOMMENT THE COMMANDS BELOW TO TEST DIFFERENT FUNCTIONALITY ***

	// TEST: CRUD operations
	// call each email type (for testing purposes)
	// newEmail := createEmail(client, "testerGal@test.ca")
	// newEmail.ConfirmedAt = 10000
	// updateEmail(client, *newEmail)
	// deleteEmail(client, newEmail.Email)
	// getEmailBatch(client, 5, 1)

	// TEST: Pagination
	// getEmailBatch(client, 3, 1)
	// getEmailBatch(client, 3, 2)
	// getEmailBatch(client, 3, 3)
}