package grpcapi

import (
	"database/sql"
	"time"

	"github.com/IM-Deane/mailing-list/mdb"
	pb "github.com/IM-Deane/mailing-list/proto"
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