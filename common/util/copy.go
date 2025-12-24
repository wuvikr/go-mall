package util

import (
	"errors"
	"go-mall/common/enum"
	"regexp"
	"time"

	"github.com/jinzhu/copier"
)

func CopyProperties(dest, src any) error {
	err := copier.CopyWithOption(dest, src, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
		Converters: []copier.TypeConverter{
			{
				// time.Time to string
				SrcType: time.Time{},
				DstType: copier.String,
				Fn: func(src any) (dst any, err error) {
					s, ok := src.(time.Time)
					if !ok {
						return nil, errors.New("src type is not time.Time")
					}
					return s.Format(enum.TimeFormatHyphenedYMDHIS), nil
				},
			},
			{
				// string to time.Time
				SrcType: copier.String,
				DstType: time.Time{},
				Fn: func(src any) (dst any, err error) {
					s, ok := src.(string)
					if !ok {
						return nil, errors.New("src type is not string")
					}

					pattern := `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$` // YYYY-MM-DD HH:MM:SS
					matched, _ := regexp.MatchString(pattern, s)
					if matched {
						return time.Parse(enum.TimeFormatHyphenedYMDHIS, s)
					}
					return nil, errors.New("src type is not time format string")
				},
			},
		},
	})

	return err
}
