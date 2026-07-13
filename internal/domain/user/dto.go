package user

type CreateUserInput struct {
	Name     string
	Email    string
	Password string
}

type UpdateUserInput struct {
	Name  *string
	Email *string
}

type UpdatePasswordInput struct {
	OldPassword string
	NewPassword string
}
