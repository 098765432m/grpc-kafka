package utils

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ParsePgDate(dateStr string) (pgtype.Date, error) {

	time, err := time.Parse("01-02-2006", dateStr)
	if err != nil {
		return pgtype.Date{}, err
	}

	var pgDate pgtype.Date
	err = pgDate.Scan(time)
	if err != nil {
		return pgtype.Date{}, err
	}

	return pgDate, nil
}

func ParseUUIDArray(uuidsStr []string) ([]pgtype.UUID, error) {

	var uuids []pgtype.UUID
	for _, uuidStr := range uuidsStr {
		var uuid pgtype.UUID
		if err := uuid.Scan(uuidStr); err != nil {
			return nil, fmt.Errorf("invalid UUID format")
		}

		uuids = append(uuids, uuid)
	}

	return uuids, nil
}
