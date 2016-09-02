package Dispatcher

// Job holds the attributes needed to perform unit of work.
type Job struct {
	Ukey   string
	HubId  int64
	NodeId int64
	CpJson []byte
}
