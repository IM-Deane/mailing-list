package grpcapi

import (
	"context"
	"database/sql"
	"log"
	"net"
	"time"

	"github.com/IM-Deane/mailing-list/mdb"
	pb "github.com/IM-Deane/mailing-list/proto"
	"google.golang.org/grpc"
)

type MailServer struct {
	pb.UnimplementedMailingListServiceServer
	db *sql.DB
}

// pbEntryToMdbEntry accepts protocol buffer and converts to mailing database EmailEntry
func pbEntryToMdbEntry(pbEntry *pb.EmailEntry) mdb.EmailEntry {
	t := time.Unix(pbEntry.ConfirmedAt, 0)
	return mdb.EmailEntry{
		ID: pbEntry.Id,
		Email: pbEntry.Email,
		ConfirmedAt: &t,
		OptOut: pbEntry.OptOut,
	}
}


// mdbEntryToPbEntry accepts an mailing database pointer and converts to a protocol buffer
func mdbEntryToPbEntry(mdbEntry *mdb.EmailEntry) *pb.EmailEntry {
	return &pb.EmailEntry{
		Id: mdbEntry.ID,
		Email: mdbEntry.Email,
		ConfirmedAt: mdbEntry.ConfirmedAt.Unix(),
		OptOut: mdbEntry.OptOut,
	}
}

// emailResponse get email, convert to protocol buffer and return
func emailResponse(db *sql.DB, email string) (*pb.EmailResponse, error) {
	entry, err := mdb.GetEmail(db, email)
	if err != nil {
		return &pb.EmailResponse{}, err
	}

	if entry == nil {
		return &pb.EmailResponse{}, nil
	}

	// convert to protocol buffer
	res := mdbEntryToPbEntry(entry)

	return &pb.EmailResponse{EmailEntry: res}, nil
}

// GetEmail gRPC handler for fetching an email
func (s *MailServer) GetEmail(ctx context.Context, req *pb.GetEmailRequest) (*pb.EmailResponse, error) {
	log.Printf("gRPC GetEmail: %v\n", req)
	return emailResponse(s.db, req.EmailAddr)
}

// GetEmailBatch gRPC handler for fetching a batch of emails
func (s *MailServer) GetEmailBatch(ctx context.Context, req *pb.GetEmailBatchRequest) (*pb.GetEmailBatchResponse, error) {
	log.Printf("gRPC GetEmailBatch: %v\n", req)

	params := mdb.GetEmailBatchQueryParams{
		Page: int(req.Page),
		Count: int(req.Count),
	}

	// query DB for emails
	mdbEntries, err := mdb.GetEmailBatch(s.db, params)
	if err != nil {
		return &pb.GetEmailBatchResponse{}, err
	}

	// create slice of email entries
	pbEntries := make([]*pb.EmailEntry, 0, len(mdbEntries))
	for i := 0; i < len(mdbEntries); i++ {
		// convert to protocol buffer entry
		entry := mdbEntryToPbEntry(&mdbEntries[i])
		// add to list
		pbEntries = append(pbEntries, entry)
	}

	return &pb.GetEmailBatchResponse{EmailEntry: pbEntries}, nil
}

// CreateEmail gRPC handler for creating an email via gRPC
func (s *MailServer) CreateEmail(ctx context.Context, req *pb.CreateEmailRequest) (*pb.EmailResponse, error) {
	log.Printf("gRPC CreateEmail: %v\n", req)

	// create new email entry in DB
	err := mdb.CreateEmail(s.db, req.EmailAddr)
	if err != nil {
		return &pb.EmailResponse{}, err
	}

	return emailResponse(s.db, req.EmailAddr)
}

// UpdateEmail gRPC handler for creating an email via gRPC
func (s *MailServer) UpdateEmail(ctx context.Context, req *pb.UpdateEmailRequest) (*pb.EmailResponse, error) {
	log.Printf("gRPC UpdateEmail: %v\n", req)

	// convert to DB entry
	entry := pbEntryToMdbEntry(req.EmailEntry)

	// update email entry in DB
	err := mdb.UpdateEmail(s.db, entry)
	if err != nil {
		return &pb.EmailResponse{}, err
	}

	return emailResponse(s.db, entry.Email)
}

// DeleteEmail gRPC handler for removing an email via gRPC
func (s *MailServer) DeleteEmail(ctx context.Context, req *pb.DeleteEmailRequest) (*pb.EmailResponse, error) {
	log.Printf("gRPC DeleteEmail: %v\n", req)

	// remove email entry in DB
	err := mdb.DeleteEmail(s.db, req.EmailAddr)
	if err != nil {
		return &pb.EmailResponse{}, err
	}

	return emailResponse(s.db, req.EmailAddr)
}

// Serve serves the gRPC handlers
func Serve(db *sql.DB, bind string) {
	// bind to address
	listener, err := net.Listen("tcp", bind)
	if err != nil {
		log.Fatalf("gRPC server error: failure to bind %v\n", bind)
	}

	// create servers
	gRPCServer := grpc.NewServer()
	mailServer := MailServer{db: db}

	// register servers
	pb.RegisterMailingListServiceServer(gRPCServer, &mailServer)

	// start server
	log.Printf("gRPC API server listening on %v\n", bind)
	if err := gRPCServer.Serve(listener); err != nil {
		log.Fatalf("gRPC server error: %v\n", err)
	}
}