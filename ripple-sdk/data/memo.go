package data

import "fmt"

const memotable = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~:/?#[]@!$&'()*+,;=%"

type Memo struct {
	Memo struct {
		MemoType   VariableLength
		MemoData   VariableLength
		MemoFormat VariableLength
	}
}

type Memos []Memo

func (m *Memo) GetTypeString() string {
	return string(m.Memo.MemoType)
}

func (m *Memo) GetDataString() string {
	return string(m.Memo.MemoData)
}

func (m *Memo) GetFormatString() string {

	return string(m.Memo.MemoFormat)
}

func (m *Memo) SetTypeFromString(s string) error {
	if m == nil {
		return fmt.Errorf("Memo is nil")
	}
	m.Memo.MemoType = []byte(s)
	return nil
}

func (m *Memo) SetDataFromString(s string) error {
	if m == nil {
		return fmt.Errorf("Memo is nil")
	}
	m.Memo.MemoData = []byte(s)

	return nil
}

func (m *Memo) SetFormatFromString(s string) error {
	if m == nil {
		return fmt.Errorf("Memo is nil")
	}
	m.Memo.MemoFormat = []byte(s)
	return nil
}
