package dbx

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//IfNotFound 判断tx是否返回 gorm.ErrRecordNotFound 错误，是的话，用repl替代
func IfNotFound(tx *gorm.DB, repl ...error) error {
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			if len(repl) > 0 {
				return fmt.Errorf("%w: %v", repl[0], tx.Error)
			}
			return nil
		}
		return tx.Error
	}
	return nil
}

func DoUpdateColumns(columns ...string) clause.OnConflict {
	return clause.OnConflict{
		DoUpdates: clause.AssignmentColumns(columns),
	}
}

var DoNothing = clause.OnConflict{DoNothing: true}

