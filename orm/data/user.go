package data

type BaseInfo struct {
	Detail      string `json:"detail"`
	Description string `json:"description"`
}

type User struct {
	Id        int                `json:"Id" orm:"u_id"` // 使用标签自定义列名
	Username  string             `json:"Username" orm:"u_username"`
	Email     string             `json:"Email" orm:"email"`
	Birthdate string             `json:"Birthdate"`
	IsActive  bool               `json:"IsActive"`
	BaseInfo  JsonData[BaseInfo] `json:"base_info"`
}

// TableName 自定义表明
func (u *User) TableName() string {
	return "users"
}
