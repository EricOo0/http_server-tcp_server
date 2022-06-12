package model

//user_info 表
type UserInfo struct {
	Id              int64  `xorm:"pk autoincr BIGINT(20)"`
	UserId          string `xorm:"not null default '' varchar(50) "`
	Password        string `xorm:"not null default '' varchar(50) "`
	Nickname        string `xorm:"not null default '' varchar(10) "`
	ProfilePicture  string `xorm:"not null default '' varchar(50) "`
	UpdateTimestamp int64  `xorm:"not null default 0 BIGINT(20)"`
	CreateTimestamp int64  `xorm:"not null default 0 BIGINT(20)"`
}
type UserInfoRecord struct {
	Id              *int64
	UserId          *string
	Password        *string
	Nickname        *string
	ProfilePicture  *string
	UpdateTimestamp *int64
	CreateTimestamp *int64
}

func GenerateUserinfoRecord(user UserInfo) UserInfoRecord {
	record := UserInfoRecord{
		UserId: &user.UserId,
	}
	if user.Password != "" {
		record.Password = &user.Password
	}
	if user.Nickname != "" {
		record.Nickname = &user.Nickname
	}
	if user.UpdateTimestamp != 0 {
		record.UpdateTimestamp = &user.UpdateTimestamp
	}
	return record
}
