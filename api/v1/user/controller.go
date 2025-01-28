package user

import (
	"bytes"
	"io"
	"path/filepath"
	"strings"

	"github.com/Npwskp/GymsbroBackend/api/v1/function"
	userFitnessPreferenceEnums "github.com/Npwskp/GymsbroBackend/api/v1/user/enums"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	Instance fiber.Router
	Service  IUserService
}

type Error error

// @Summary		Create a user
// @Description	Create a user
// @Tags		users
// @Accept		json
// @Produce		json
// @Param		user body CreateUserDto true "Create User"
// @Success		201	{object} User
// @Failure		400	{object} Error
// @Router		/user [post]
func (uc *UserController) PostUsersHandler(c *fiber.Ctx) error {
	validate := validator.New()
	user := new(CreateUserDto)
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validate.Struct(*user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	res, err := uc.Service.CreateUser(user)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

// @Summary		Get all users
// @Description	Get all users
// @Tags		users
// @Accept		json
// @Produce		json
// @Success		200	{array}	User
// @Failure		400	{object} Error
// @Router		/user [get]
func (uc *UserController) GetAllUsersHandler(c *fiber.Ctx) error {
	users, err := uc.Service.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(users)
}

// @Summary		Get a user
// @Description	Get a user
// @Tags		users
// @Accept		json
// @Produce		json
// @Success		200	{object} User
// @Failure		400	{object} Error
// @Router		/user/me [get]
func (uc *UserController) GetUserHandler(c *fiber.Ctx) error {
	id := function.GetUserIDFromContext(c)
	user, err := uc.Service.GetUser(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

// @Summary		Get a user energy consume plan
// @Description	Get a user energy consume plan
// @Tags		users
// @Accept		json
// @Produce		json
// @Success		200	{object} userFitnessPreferenceEnums.EnergyConsumptionPlan
// @Failure		400	{object} Error
// @Router		/user/energyplan [get]
func (uc *UserController) GetUserEnergyConsumePlanHandler(c *fiber.Ctx) error {
	id := function.GetUserIDFromContext(c)
	plan, err := uc.Service.GetUserEnergyConsumePlan(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(plan)
}

// @Summary		Get all activity levels
// @Description	Get all activity levels
// @Tags		users
// @Accept		json
// @Produce		json
// @Success		200	{array}	userFitnessPreferenceEnums.ActivityLevelType
// @Failure		400	{object} Error
// @Router		/user/activitylevels [get]
func (uc *UserController) GetAllActivityLevels(c *fiber.Ctx) error {
	levels := userFitnessPreferenceEnums.GetAllActivityLevels()
	return c.Status(fiber.StatusOK).JSON(levels)
}

// @Summary		Get all goals
// @Description	Get all goals
// @Tags		users
// @Accept		json
// @Produce		json
// @Success		200	{array}	userFitnessPreferenceEnums.GoalType
// @Failure		400	{object} Error
// @Router		/user/goals [get]
func (uc *UserController) GetAllGoals(c *fiber.Ctx) error {
	goals := userFitnessPreferenceEnums.GetAllGoals()
	return c.Status(fiber.StatusOK).JSON(goals)
}

// @Summary		Get all carb preferences
// @Description	Get all carb preferences
// @Tags		users
// @Accept		json
// @Produce		json
// @Success		200	{array}	userFitnessPreferenceEnums.CarbPreferenceType
// @Failure		400	{object} Error
// @Router		/user/carbpreferences [get]
func (uc *UserController) GetAllCarbPreferences(c *fiber.Ctx) error {
	carbPreferences := userFitnessPreferenceEnums.GetAllCarbPreferences()
	return c.Status(fiber.StatusOK).JSON(carbPreferences)
}

// @Summary		Delete a user
// @Description	Delete a user
// @Tags		users
// @Accept		json
// @Produce		json
// @Success		204
// @Failure		400	{object} Error
// @Router		/user/me [delete]
func (uc *UserController) DeleteUserHandler(c *fiber.Ctx) error {
	id := function.GetUserIDFromContext(c)
	err := uc.Service.DeleteUser(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary		Update a user username and password
// @Description	Update a user username and password
// @Tags		users
// @Accept		json
// @Produce		json
// @Param		user body UpdateUsernamePasswordDto true "UpdateUsernamePassword User"
// @Success		200	{object} User
// @Failure		400	{object} Error
// @Router		/user/usepass [patch]
func (uc *UserController) UpdateUsernamePassword(c *fiber.Ctx) error {
	id := function.GetUserIDFromContext(c)
	validate := validator.New()
	doc := new(UpdateUsernamePasswordDto)
	if err := c.BodyParser(&doc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validate.Struct(*doc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	user, err := uc.Service.UpdateUsernamePassword(doc, id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

// @Summary		Update a user body
// @Description	Update a user body
// @Tags		users
// @Accept		json
// @Produce		json
// @Param		user body UpdateBodyDto true "UpdateBody User"
// @Success		200	{object} User
// @Failure		400	{object} Error
// @Router		/user/body [patch]
func (uc *UserController) UpdateBody(c *fiber.Ctx) error {
	id := function.GetUserIDFromContext(c)
	validate := validator.New()
	doc := new(UpdateBodyDto)
	if err := c.BodyParser(&doc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := validate.Struct(*doc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	user, err := uc.Service.UpdateBody(doc, id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

// @Summary     Update first login status
// @Description Mark user as not first time login
// @Tags        users
// @Accept      json
// @Produce     json
// @Success     204
// @Failure     400 {object} Error
// @Router      /user/first-login [patch]
func (uc *UserController) UpdateFirstLoginStatus(c *fiber.Ctx) error {
	id := function.GetUserIDFromContext(c)
	err := uc.Service.UpdateFirstLoginStatus(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary     Update user profile picture
// @Description Upload and update user's profile picture
// @Tags        users
// @Accept      multipart/form-data
// @Produce     json
// @Param       file formData file true "Profile picture (jpeg/png)"
// @Success     204
// @Failure     400 {object} Error
// @Router      /user/picture [patch]
func (uc *UserController) UpdateUserPicture(c *fiber.Ctx) error {
	id := function.GetUserIDFromContext(c)

	// Get the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No file uploaded",
		})
	}

	// Validate file type (optional but recommended)
	contentType := file.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Only JPEG and PNG files are allowed",
		})
	}

	// Validate file size (e.g., max 5MB)
	if file.Size > 5*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "File size exceeds 5MB limit",
		})
	}

	// Open and read the file
	fileHandle, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to process uploaded file",
		})
	}
	defer fileHandle.Close()

	// Read file contents
	fileBytes, err := io.ReadAll(fileHandle)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to read file contents",
		})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExtensions := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
	}

	contentType, isAllowed := allowedExtensions[ext]
	if !isAllowed {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Unsupported image format. Allowed formats: JPG, PNG, GIF, WEBP",
		})
	}

	// Update the user picture
	user, err := uc.Service.UpdateUserPicture(c, id, bytes.NewReader(fileBytes), file.Filename, contentType)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (uc *UserController) Handle() {
	g := uc.Instance.Group("/user")

	g.Post("", uc.PostUsersHandler)
	g.Get("", uc.GetAllUsersHandler)
	g.Get("/me", uc.GetUserHandler)
	g.Get("/energyplan", uc.GetUserEnergyConsumePlanHandler)
	g.Get("/activitylevels", uc.GetAllActivityLevels)
	g.Get("/goals", uc.GetAllGoals)
	g.Get("/carbpreferences", uc.GetAllCarbPreferences)
	g.Delete("/me", uc.DeleteUserHandler)
	g.Patch("/body", uc.UpdateBody)
	g.Patch("/usepass", uc.UpdateUsernamePassword)
	g.Patch("/first-login", uc.UpdateFirstLoginStatus)
	g.Patch("/picture", uc.UpdateUserPicture)
}
