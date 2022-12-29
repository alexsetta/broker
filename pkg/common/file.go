package common

type File struct {
	Config string
	Wallet string
	Last   string
	Quotes string
	Order  string
	Log    string
}

func NewFile(dir *Dir) *File {
	return &File{
		Config: dir.Config + "broker.cfg",
		Wallet: dir.Config + "wallet.cfg",
		Last:   dir.Files + "last.txt",
		Quotes: dir.Files + "quotes.log",
		Order:  dir.Files + "order.log",
		Log:    dir.Files + "broker.log",
	}
}
