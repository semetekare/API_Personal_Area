package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(100)" json:"name"`
	NdsLogin  string    `gorm:"type:varchar(128)" json:"nds_login"`
	RolID     int       `json:"rol_id"`
	UAccess   int       `json:"u_access"` // Используем int вместо tinyint
	FacID     int       `json:"fac_id"`
	KafID     int       `json:"kaf_id"`
	SpecID    int       `json:"spec_id"`
	DepID     int       `gorm:"column:dep_id" json:"dep_id"` // Поле с аннотацией для соответствия столбцу с другим именем
	DogdepID  int       `gorm:"column:dogdep_id" json:"dogdep_id"`
	Password  string    `gorm:"type:varchar(100)" json:"-"`
	UType     int       `gorm:"column:u_type" json:"u_type"`
	QRToken   string    `gorm:"type:text" json:"-"`
	PersonID  int       `gorm:"column:person_id" json:"person_id"`
	FacList   string    `gorm:"type:varchar(255)" json:"fac_list"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Метод для автоматического обновления времени при сохранении
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()
	return
}

type Employee struct {
	ID           int        `gorm:"primary_key;auto_increment" json:"id"`
	Name1        string     `json:"name1"`
	Name2        string     `json:"name2"`
	Name3        string     `json:"name3"`
	NDSLogin     string     `json:"nds_login"`
	PersonID     int        `json:"person_id"`
	EmpDolID     int        `json:"emp_dol_id"` // 1 - преподаватель, 2 - зав. каф, 3 - сотрудник уму
	KafID        int        `json:"kaf_id"`
	FIO          string     `json:"fio"`
	ADLogin      bool       `json:"ad_login"` // Флаг наличия административного логина
	UMUIDOld     int        `json:"umu_id_old"`
	UMUID        int        `json:"umu_id"`
	BirLoc       string     `json:"bir_loc"`
	Comments     string     `json:"comments"`
	DipDate      *time.Time `json:"dip_date"`
	DipNum       string     `json:"dip_num"`
	DOB          *time.Time `json:"dob"`
	Family       string     `json:"family"`
	FStatID      int        `json:"fstat_id"`
	KladrString  string     `json:"kladr_string"`
	KvalifyID    int        `json:"kvalify_id"`
	OutDoc       *time.Time `json:"out_doc"`
	PassDate     *time.Time `json:"pass_date"`
	PassLoc      string     `json:"pass_loc"` // Место выдачи паспорта
	Passport     string     `json:"passport"`
	Phone        string     `json:"phone"`
	Staj         int        `json:"staj"`
	UZID         int        `json:"uz_id"`
	PassportS    string     `json:"passport_s"`  // Серия паспорта
	Sex          int        `json:"sex"`         // Пол
	StajIRGUPS   int        `json:"staj_irgups"` // Стаж в ИрГУПС
	SecretAPIKey string     `json:"secret_api_key"`
	QRToken      string     `json:"qr_token"`  // Секретный ключ для QR
	Metaphone    string     `json:"metaphone"` // Фонетическое звучание фамилии
	UpdatedAt    *time.Time `json:"updated_at"`
	CreatedAt    time.Time  `json:"created_at"`
	Base1C       int        `json:"base_1c"`  // База 1С
	Password     string     `json:"password"` // Временный пароль для теста
}
