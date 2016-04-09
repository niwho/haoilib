package haoilib

type HB interface {
	SetRebateW(softkey string)
	SendByteExW(MyUserStr, GameID string, PicBuffer []byte, Size, TimeOut, LostPoint int64,
		BeiZhu string) (Result, Reply string, err error)
	GetPointW(MyUserStr string) (int, error)
	GetBusyW() (string, error)
}
