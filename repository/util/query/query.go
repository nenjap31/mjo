package query

import (
	"fmt"
	"net/url"
	"strings"

	//"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gorm"
)

func GormFilterBuilder(db *gorm.DB, filters url.Values, limit int, offeset int) *gorm.DB {
	for key, values := range filters {
		for _, value := range values {
			keySplitter := strings.Split(key, "__")
			field := keySplitter[0]
			operator := keySplitter[1]
			switch operator {
			case "in":
				db = db.Where(fmt.Sprintf("`%s` IN ?", field), values)
			case "not_in":
				db = db.Where(fmt.Sprintf("`%s` NOT IN ?", field), values)
			case "equals":
				db = db.Where(fmt.Sprintf("`%s` = ?", field), value)
			case "not_equals":
				db = db.Where(fmt.Sprintf("`%s` != ?", field), value)
			case "gt":
				db = db.Where(fmt.Sprintf("`%s` > ?", field), value)
			case "gte":
				db = db.Where(fmt.Sprintf("`%s` >= ?", field), value)
			case "lt":
				db = db.Where(fmt.Sprintf("`%s` < ?", field), value)
			case "lte":
				db = db.Where(fmt.Sprintf("`%s` <= ?", field), value)
			case "is_null":
				db = db.Where(fmt.Sprintf("`%s` IS NULL", field))
			case "is_not_null":
				db = db.Where(fmt.Sprintf("`%s` IS NOT NULL", field))
			}
			continue
		}
	}
	return db.Limit(limit).Offset(offeset)
}
