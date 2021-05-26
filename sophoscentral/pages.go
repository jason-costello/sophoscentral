package sophoscentral

type PageByOffsetOptions struct {
ListByPageOffset
}

type PageByOffset struct {
Current *int `json:"current,omitempty"`
Total   *int `json:"total,omitempty"`
Size    *int `json:"size,omitempty"`
Maxsize *int `json:"maxSize,omitempty"`
Items   *int `json:"items,omitempty"`
}
