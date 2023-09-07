package Handlers

import (
	"api_sotr/app/service"

	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	EName string `json:"e_name"`
	EPass string `json:"e_pass"`
}

type LoginResponse struct {
	UserID     int    `json:"user_id"`
	UserName   string `json:"user_name"`
	UserAccess string `json:"user_access"`
}

func Login_Employee(c *fiber.Ctx) error {
	request := new(LoginRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Некорректные данные"})
	}

	evaName := request.EName
	evaPass := request.EPass

	if len(evaName) > 3 && len(evaPass) > 3 {

		checkUserLogin := service.CheckUserLogin(evaName, evaPass)
		if checkUserLogin {
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

		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Некорректный пароль / имя пользователя"})
		}

	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Короткий пароль или имя пользователя"})
	}
	return nil
}
