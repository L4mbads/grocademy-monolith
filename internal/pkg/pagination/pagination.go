package pagination

import (
	"fmt"
	"math"

	"gorm.io/gorm"
)

type Pagination struct {
	CurrentPage int64 `json:"current_page"`
	TotalPages  int64 `json:"total_pages"`
	TotalItems  int64 `json:"total_items"`
}

func Paginate(db *gorm.DB, dest any, page, limit int64, searchableColumns []string, query string) (*any, Pagination, error) {
	var totalItems int64

	countDB := db
	if query != "" && len(searchableColumns) > 0 {
		searchQuery := ""
		for i, col := range searchableColumns {
			searchQuery += fmt.Sprintf("%s ILIKE ?", col) // ILIKE is case-insensitive
			if i < len(searchableColumns)-1 {
				searchQuery += " OR "
			}
		}

		args := make([]interface{}, len(searchableColumns))
		for i := range searchableColumns {
			args[i] = fmt.Sprintf("%%%s%%", query)
		}

		db = db.Where(searchQuery, args...)
		countDB = countDB.Where(searchQuery, args...)
	}

	pagination := Pagination{
		CurrentPage: 0,
		TotalPages:  0,
		TotalItems:  0,
	}

	if err := countDB.Model(dest).Count(&totalItems).Error; err != nil {
		return nil, pagination, err
	}

	totalPages := int64(math.Ceil(float64(totalItems) / float64(limit)))
	if totalPages == 0 && totalItems > 0 {
		totalPages = 1
	}

	if page < 1 {
		page = 1
	}
	if totalPages > 0 && page > totalPages {
		page = totalPages
	}

	pagination.CurrentPage = page
	pagination.TotalItems = totalItems
	pagination.TotalPages = totalPages

	offset := (page - 1) * limit

	// Execute the query with limit and offset
	if err := db.Limit(int(limit)).Offset(int(offset)).Find(dest).Error; err != nil {
		return nil, pagination, err
	}

	return &dest, pagination, nil
}
