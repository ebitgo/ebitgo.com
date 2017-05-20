package dbmodels

import (
	"time"

	"github.com/ebitgo/ebitgo.com/models/logger"
	"github.com/go-xorm/xorm"
)

type UserIdT struct {
	Id    uint64 `xorm:"pk"`
	Email string `xorm:"unique"`
}

type UserDbT struct {
	Id            uint64 `xorm:"pk"`
	Password      string
	MobilePhone   string
	UseGa         bool
	GaKey         string
	RegTime       time.Time `xorm:"DateTime"`
	LastLoginTime time.Time `xorm:"DateTime"`
	AuthStr       string
	UserLevel     uint64
	ActiveFlag    string
}

type WalletDbT struct {
	Id         uint64 `xorm:"pk"`
	UserId     uint64
	NickName   string
	UseGa      bool
	PublicAddr string
	SecretKey  string
	RegTime    time.Time `xorm:"DateTime"`
}

type FriendlyShipDbT struct {
	UserId     uint64
	WalletId   uint64
	NickName   string
	PublicAddr string
}

type FriendlyGroupDbT struct {
	UserId    uint64
	GroupId   uint64
	WalletId  uint64
	GroupName string
}

type WalletMessageDbT struct {
	MsgId          string `xorm:"pk"`
	FromUserId     uint64
	FromWalletId   uint64
	FromWalletAddr string
	ToWalletId     uint64
	ToWalletAddr   string
	Agent          bool
	CreateTime     time.Time `xorm:"DateTime"`
	SendTime       time.Time `xorm:"DateTime"`
	Context        string
	IsRead         bool
	IsSuccess      bool
}

func OrmRegModels(eng *xorm.Engine) {
	err := eng.Sync(new(UserIdT), new(UserDbT), new(WalletDbT), new(FriendlyShipDbT), new(FriendlyGroupDbT))
	if err != nil {
		logmodels.Logger.Info(logmodels.SPrintError("Database Table", "XORM Engine Sync is err %v", err))
		panic(1)
	}
}
