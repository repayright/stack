// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

type Cursor struct {
	Data     []PaymentsAccount `json:"data"`
	HasMore  bool              `json:"hasMore"`
	Next     *string           `json:"next,omitempty"`
	PageSize int64             `json:"pageSize"`
	Previous *string           `json:"previous,omitempty"`
}

func (o *Cursor) GetData() []PaymentsAccount {
	if o == nil {
		return []PaymentsAccount{}
	}
	return o.Data
}

func (o *Cursor) GetHasMore() bool {
	if o == nil {
		return false
	}
	return o.HasMore
}

func (o *Cursor) GetNext() *string {
	if o == nil {
		return nil
	}
	return o.Next
}

func (o *Cursor) GetPageSize() int64 {
	if o == nil {
		return 0
	}
	return o.PageSize
}

func (o *Cursor) GetPrevious() *string {
	if o == nil {
		return nil
	}
	return o.Previous
}

type AccountsCursor struct {
	Cursor Cursor `json:"cursor"`
}

func (o *AccountsCursor) GetCursor() Cursor {
	if o == nil {
		return Cursor{}
	}
	return o.Cursor
}