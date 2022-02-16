package store

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
)

type dbservice struct {
	client *FirestoreClient
}

type FirestoreClient struct {
	fs *firestore.Client

	common                   dbservice
	StandingsService         *StandingsService

}

func NewClient(ctx context.Context, projectId string) (*FirestoreClient, error) {

	firestore, err := firestore.NewClient(ctx, projectId)

	if err != nil {
		return nil, err
	}

	fsc := &FirestoreClient{
		fs: firestore,
	}

	fsc.common.client = fsc
	fsc.StandingsService = (*StandingsService)(&fsc.common)
	return fsc, nil
}

func parseBool(val string) bool {
	flag, err := strconv.ParseBool(val)
	if err != nil {
		flag = false
	}
	return flag
}

func parseInt(val string) int {
	result, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return result
}

func parseInt64(val string) int64 {
	result, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0
	}
	return result
}

func parseDate(format string, val string) time.Time {
	mydate, err := time.Parse(format, val)
	if err != nil {
		return time.Now().UTC()
	}
	return mydate.UTC()
}

func DocWithIDAndName(id int, name string) string {
	return fmt.Sprintf("%d#%s", id, strings.ToUpper(name))
}