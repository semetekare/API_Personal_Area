package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"strings"

	//"gorm.io/driver/mysql"
	"api_sotr/app/models"
)

func GetUvalStatesForEmpId(empID int) (int, error) {
	db := models.DbSotr

	var result int
	err := db.Raw(`
		SELECT e.id
		FROM employes e
		JOIN emp_dolzn ed ON ed.empolyes_id=e.id
		WHERE e.id=? AND ed.status != 0 AND ed.uval_date IS NULL
	`, empID).Scan(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("No result found")
		}
		return 0, err
	}

	return result, nil
}

// UserInfo предоставляет информацию являеться ли пользователь сотрудником.
func EmployeeInfo(ndsLogin string) (bool, int, error) {
	db := models.DbSotr

	var userCount int64
	err := db.Table("users").
		Where("nds_login = ? AND rol_id IS NOT NULL", ndsLogin).
		Count(&userCount).Error
	if err != nil {
		return false, 0, err
	}

	var employee models.Employee
	err = db.Raw(`
        SELECT emp.id, emp.name1, emp.name2, emp.name3, emp.umu_id, emp.emp_dol_id, emp.kaf_id, k.name kname, emp.person_id, f.filial_id, p.fil_id, u.id AS user_id
        FROM employes emp
        LEFT JOIN kafedras k ON k.id=emp.kaf_id
        LEFT JOIN faculties f ON k.fac_id = f.id
        LEFT JOIN persons p ON p.id = emp.person_id
        LEFT JOIN users u ON u.nds_login = emp.nds_login AND u.u_access = 1
        WHERE emp.nds_login = ?
    `, ndsLogin).Scan(&employee).Error

	if err != nil {
		return false, 0, err
	}

	return userCount > 0, employee.ID, nil
}

func GetStudentInfo(studentID, mode int) (*models.StudentInfo, error) {
	db := models.DbSotr
	var studentInfo models.StudentInfo

	if studentID <= 0 {
		return nil, fmt.Errorf("ERROR: Неверно задан входной параметр!")
	}

	if mode == 0 {
		query := `
			SELECT s.id AS stud_id, s.name1 AS fname, s.name2 AS lname, s.name3 AS mname, s.nzk, 
			DATE_FORMAT(s.dob, '%d-%m-%Y') AS dob, s.person_id, p.phone, p.country, sc.nds_login
			FROM students s
			JOIN persons p ON p.id = s.person_id
			JOIN stud_cards sc ON sc.id = s.id
			WHERE s.id = ?`

		if err := db.Raw(query, studentID).Scan(&studentInfo).Error; err != nil {
			return nil, err
		}
	} else if mode == 1 {
		query := `
			SELECT sg.name AS gname, f.short_name AS facname, sf.name AS sfname, sg.id AS sgid, 
			f.id AS facid, sf.id AS sfid, p.name1, p.name2, p.name3, sp.name AS spname, 
			sp1.name AS spparent, sp.kod_spec AS shifr, f.cheef_dolzn AS cheefdolzn, 
			f.cheef_name AS cheefname, c.course, ef.name AS eduform, sg.is_filial, 
			sc.person_id, c.archived_id AS sstatus, sg.archived_id AS gstatus
			FROM stud_cards sc
			JOIN persons p ON sc.person_id = p.id
			JOIN courses c ON c.stud_id = sc.id
			JOIN stud_groups sg ON sg.id = c.group_id
			JOIN faculties f ON f.id = sg.fac_id
			JOIN studyforms sf ON sf.id = sg.studyform_id
			JOIN specialities sp ON sp.id = sg.spec_id
			JOIN specialities sp1 ON sp1.id = sp.parent_id
			JOIN edu_form ef ON ef.id = sg.edu_form_id
			WHERE c.active = 1 AND c.stud_id = ?`

		if err := db.Raw(query, studentID).Scan(&studentInfo).Error; err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("ERROR: Неверное значение параметра 'mode'")
	}

	return &studentInfo, nil
}

// проверка студента, с возращаемой маппой(id, person_id)
func AsuIsStudent(ndsLogin string) (interface{}, error) {
	db := models.DbSotr

	ndsLogin = strings.TrimSpace(ndsLogin)
	if len(ndsLogin) < 5 {
		return nil, errors.New("ERROR: Слишком короткий логин")
	}

	var studentData struct {
		ID       int
		PersonID int
	}

	err := db.Raw("SELECT id, person_id FROM stud_cards WHERE nds_login = ?", ndsLogin).Scan(&studentData).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Нет записей, возвращаем nil
		}
		return nil, err
	}

	return map[string]interface{}{"id": studentData.ID, "person_id": studentData.PersonID}, nil
}

// checkMySQLAuth проверяет аутентификацию сотрудника через MySQL.
func checkMySQLAuth(db *gorm.DB, username, password string) (bool, error) {
	var user models.User
	err := db.Raw("SELECT id FROM users WHERE password=md5(?) AND nds_login=? AND u_access=?", password, username, 1).Scan(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil // Нет записей, доступ запрещен
		}
		return false, err
	}
	return true, nil
}

// checkOtherCase проверяет аутентификацию пользователя студент
func checkOtherCase(db *gorm.DB, username, password string) (bool, error) {
	var pass string
	err := db.Raw(`
		SELECT CONCAT('17March_', sc.id) as pass
		FROM stud_cards sc
		JOIN courses c ON sc.id = c.stud_id
		WHERE nds_login = ? AND c.active = 1 AND
			(c.archived_id != 1 AND c.archived_id != 2 AND c.archived_id != 4)
	`, username).Scan(&pass).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil // Нет записей, доступ запрещен
		}
		return false, err
	}

	if password == pass {
		return true, nil
	}

	return false, nil
}

/*
*	CheckUserLogin выполняет аутентификацию пользователя.
*	добавить реализацию выбора проверки
 */
func CheckUserLogin(evaName, evaPass string) bool {

	access := false

	models.ConnectDatabase() // Подключение к базе данных
	db := models.DbSotr      //DbSotr - это глобальная переменная типа gorm.DB

	if evaName != "" && evaPass != "" {
		// Проверка аутентификации через MySQL
		mysqlAccess, err := checkMySQLAuth(db, evaName, evaPass)
		if err != nil {
			log.Fatal(err)
		}

		mysqlAccessStud, err := checkOtherCase(db, evaName, evaPass)
		if err != nil {
			log.Fatal(err)
		}

		// Определение общего доступа
		if (mysqlAccess == true) || (mysqlAccessStud == true) {
			access = true
			fmt.Println("Аутентификация успешна")
		} else {
			access = false
			fmt.Println("Аутентификация не удалась")
		}

		return access
	} else {
		return access
	}
}
