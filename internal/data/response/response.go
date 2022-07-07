package response

// Pagination 分页结构
type Pagination struct {
    Page      uint          `json:"page"`
    TotalPage uint          `json:"totalPage"`
    Size      uint          `json:"size"`
    Items     []interface{} `json:"items"`
    TotalSize uint          `json:"totalSize"`
}
