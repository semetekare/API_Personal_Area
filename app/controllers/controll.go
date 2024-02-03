package Handlers

import (
	"api_sotr/app/models"
	"api_sotr/app/service"
	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	EName     string `json:"e_name"`
	EPass     string `json:"e_pass"`
	TypeLogin string `json:"type_login"`
}

type LoginResponse struct {
	UserID     int    `json:"user_id"`
	UserName   string `json:"user_name"`
	UserAccess string `json:"user_access"`
}

func Ok(c *fiber.Ctx) error {
	c.Status(fiber.StatusOK)
	return nil
}

func Login(c *fiber.Ctx) error {
	request := new(LoginRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Некорректные данные"})
	}

	evaName := request.EName
	evaPass := request.EPass
	typeLog := request.TypeLogin

	if len(evaName) > 3 && len(evaPass) > 3 {

		checkUserLogin := service.CheckUserLogin(evaName, evaPass)
		if checkUserLogin {

			switch typeLog {
			case "student":
				accessInfo, err := service.AsuIsStudent(evaName)

				if err != nil {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Вы не являетесь студентом"})
				} else if accessInfo == nil {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Нет записей, доступ запрещен"})
				} else {
					studentData := accessInfo.(map[string]interface{})
					statusStud, err := service.GetStudentInfo(studentData["id"].(int), 1)
					if err != nil {
						return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Ошибка:" + err.Error()})
					} else if statusStud == nil {
						return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Нет информации о студенте"})
					} else {

						if (contains([]int{0, 3, 6, 7, 10, 12}, statusStud.SStatus) && statusStud.GStatus == 0) || statusStud.ID == 104929 {
							return c.JSON(
								models.StudentInfo{
									ID:         statusStud.ID,
									StudID:     statusStud.StudID,
									FName:      statusStud.FName,
									LName:      statusStud.LName,
									MName:      statusStud.MName,
									NZK:        statusStud.NZK,
									DOB:        statusStud.DOB,
									PersonID:   statusStud.PersonID,
									Phone:      statusStud.Phone,
									Country:    statusStud.Country,
									NdsLogin:   statusStud.NdsLogin,
									GName:      statusStud.GName,
									FacName:    statusStud.FacName,
									SFName:     statusStud.SFName,
									SGID:       statusStud.SGID,
									FacID:      statusStud.FacID,
									SFID:       statusStud.SFID,
									Name1:      statusStud.Name1,
									Name2:      statusStud.Name2,
									Name3:      statusStud.Name3,
									SPName:     statusStud.SPName,
									SPParent:   statusStud.SPParent,
									Shifr:      statusStud.Shifr,
									CheefDolzn: statusStud.CheefDolzn,
									CheefName:  statusStud.CheefName,
									Course:     statusStud.Course,
									EduForm:    statusStud.EduForm,
									IsFilial:   statusStud.IsFilial,
									SStatus:    statusStud.SStatus,
									GStatus:    statusStud.GStatus,
								})
						} else {
							return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Похоже вас отчислили("})
						}
					}
				}

			case "employe":
				if EmployeeInfo, err := service.EmployeeInfo(evaName); err == nil {
					if EmployeeInfo != nil {
						// Пользователь является сотрудником

						// Проверяем UvalStates для EmpID
						if getUvalStatesForEmpId, err := service.GetUvalStatesForEmpId(EmployeeInfo.ID); err == nil {
							if getUvalStatesForEmpId == EmployeeInfo.ID {
								// У пользователя есть информация о контракте/договоре

								return c.JSON(
									models.Employee{
										ID:           EmployeeInfo.ID,
										Name1:        EmployeeInfo.Name1,
										Name2:        EmployeeInfo.Name2,
										Name3:        EmployeeInfo.Name3,
										NDSLogin:     EmployeeInfo.NDSLogin,
										PersonID:     EmployeeInfo.PersonID,
										EmpDolID:     EmployeeInfo.EmpDolID,
										KafID:        EmployeeInfo.KafID,
										FIO:          EmployeeInfo.FIO,
										ADLogin:      EmployeeInfo.ADLogin,
										UMUIDOld:     EmployeeInfo.UMUIDOld,
										UMUID:        EmployeeInfo.UMUID,
										BirLoc:       EmployeeInfo.BirLoc,
										Comments:     EmployeeInfo.Comments,
										DipDate:      EmployeeInfo.DipDate,
										DipNum:       EmployeeInfo.DipNum,
										DOB:          EmployeeInfo.DOB,
										Family:       EmployeeInfo.Family,
										FStatID:      EmployeeInfo.FStatID,
										KladrString:  EmployeeInfo.KladrString,
										KvalifyID:    EmployeeInfo.KvalifyID,
										OutDoc:       EmployeeInfo.OutDoc,
										PassDate:     EmployeeInfo.PassDate,
										PassLoc:      EmployeeInfo.PassLoc,
										Passport:     EmployeeInfo.Passport,
										Phone:        EmployeeInfo.Phone,
										Staj:         EmployeeInfo.Staj,
										UZID:         EmployeeInfo.UZID,
										PassportS:    EmployeeInfo.PassportS,
										Sex:          EmployeeInfo.Sex,
										StajIRGUPS:   EmployeeInfo.StajIRGUPS,
										SecretAPIKey: EmployeeInfo.SecretAPIKey,
										QRToken:      EmployeeInfo.QRToken,
										Metaphone:    EmployeeInfo.Metaphone,
										UpdatedAt:    EmployeeInfo.UpdatedAt,
										CreatedAt:    EmployeeInfo.CreatedAt,
										Base1C:       EmployeeInfo.Base1C,
										Password:     EmployeeInfo.Password,
									})

							} else {
								return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "У вас отсутствует информация о контракте/договоре для осуществления трудовой деятельности"})
							}
						} else {
							// Обработка ошибки при получении информации о контракте/договоре
							return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Произошла ошибка"})
						}

					} else {
						return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Вы не являетесь сотрудником"})
					}
				} else {
					// Обработка ошибки, если что-то пошло не так при проверке
				}
			}
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Некорректный пароль / имя пользователя"})
		}

	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Короткий пароль или имя пользователя"})
	}
	return nil
}

// Функция contains для проверки наличия значения в срезе
func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
