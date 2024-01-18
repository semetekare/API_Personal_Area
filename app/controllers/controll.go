package Handlers

import (
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
							return c.JSON(LoginResponse{
								UserID:     int(statusStud.PersonID), //empID здесь
								UserName:   evaName,
								UserAccess: "allowed",
							})
						} else {
							return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Похоже вас отчислили("})
						}
					}
				}

			case "employe":
				if EmployeeInfo, empID, err := service.EmployeeInfo(evaName); err == nil {
					if EmployeeInfo {
						// Пользователь является сотрудником

						// Проверяем UvalStates для EmpID
						if getUvalStatesForEmpId, err := service.GetUvalStatesForEmpId(empID); err == nil {
							if getUvalStatesForEmpId == empID {
								// У пользователя есть информация о контракте/договоре

								return c.JSON(LoginResponse{
									UserID:     empID, //empID здесь
									UserName:   evaName,
									UserAccess: "allowed",
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
