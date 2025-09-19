package utils

import (
	"fmt"
	"time"

	"github.com/098765432m/grpc-kafka/common/consts"
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

func ToPgDateRange(dateStartStr string, dateEndStr string) (pgtype.Date, pgtype.Date, error) {
	InvalidDate := pgtype.Date{
		Valid: false,
	}

	tempTime, err := time.Parse(consts.DATE_FORMAT, dateStartStr)
	if err != nil {
		return InvalidDate, InvalidDate, err
	}

	var pgDateStart pgtype.Date
	err = pgDateStart.Scan(tempTime)
	if err != nil {
		return InvalidDate, InvalidDate, err
	}

	tempTime, err = time.Parse(consts.DATE_FORMAT, dateEndStr)
	if err != nil {
		return InvalidDate, InvalidDate, err
	}

	var pgDateEnd pgtype.Date
	err = pgDateEnd.Scan(tempTime)
	if err != nil {
		return InvalidDate, InvalidDate, err
	}

	if pgDateStart.Time.After(pgDateEnd.Time) {
		return InvalidDate, InvalidDate, fmt.Errorf("time Start is after time End")
	}

	return pgDateStart, pgDateEnd, nil

}

func ParsePgText(str string) (pgtype.Text, error) {
	var tempStr pgtype.Text
	if str == "" {
		tempStr = pgtype.Text{
			Valid: false,
		}
	} else {
		if err := tempStr.Scan(str); err != nil {
			return pgtype.Text{
				Valid: false,
			}, err
		}
	}

	return tempStr, nil
}

func ToPgInt4(number int) pgtype.Int4 {
	var tempNumber pgtype.Int4
	// Int4 scan only take int64 number, otherwise it return error
	if err := tempNumber.Scan(int64(number)); err != nil {
		return pgtype.Int4{
			Valid: false,
		}
	}

	return tempNumber
}

func ToPgUuidArray(uuidsStr []string) ([]pgtype.UUID, error) {

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

func ToPgUuidString(uuids []pgtype.UUID) []string {
	uuidsStr := make([]string, 0, len(uuids))

	for _, uuid := range uuids {
		uuidsStr = append(uuidsStr, uuid.String())
	}

	return uuidsStr
}
