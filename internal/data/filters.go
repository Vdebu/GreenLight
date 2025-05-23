package data

import (
	"greenlight.vdebu.net/internal/validator"
	"math"
	"strings"
)

// 存储其他请求端点也可能用上的字段信息
type Filters struct {
	Page         int
	PageSize     int
	Sort         string   // 存储输入排序规则
	SortSafeList []string // 存储允许的排序规则
}

// 存储当前页面的数据信息
type MetaData struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

// 检查filters字段的有效性
func ValidateFilters(v *validator.Validator, f Filters) {
	// 检查 pages
	v.Check(f.Page > 0, "page", "must greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "pages_size", "must greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
	// 检查sort的key是否都在允许的范围内
	v.Check(validator.In(f.Sort, f.SortSafeList...), "sort", "invalid sort value")
}

// 检查用户提供的sort key是可用的并将关键字进行提取
func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafeList {
		if f.Sort == safeValue {
			// 去除前缀返回有效的字符串
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	panic("unsafe sort parameter:" + f.Sort)
}

// 检测是升序还是降序排序
func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		// "-"表示降序
		return "DESC"
	}
	// 升序
	return "ASC"
}

// 获取每页显示的数据数目
func (f Filters) limit() int {
	return f.PageSize
}

// 获取当前要展示的页
func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

// 计算页面的数据信息
func calculateMetadata(totalRecords, page, pageSize int) MetaData {
	if totalRecords == 0 {
		// 如果总记录数为零返回空的结构体
		return MetaData{}
	}
	return MetaData{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))), // 计算以当前规格要显示多少页 最后的结果向上取整
		TotalRecords: totalRecords,
	}
}
