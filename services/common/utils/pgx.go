package utils

import (
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
