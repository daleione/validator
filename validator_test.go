package validator

import (
	"testing"
)

type User struct {
	Name  string
	Email string
	Age   int
}

func TestStruct(t *testing.T) {
	user := User{
		Name:  "",
		Email: "invalid_email",
		Age:   25,
	}

	structValidator := &StructValidator{}
	structValidator.AddField("Name", user.Name, Required(), MinLength(3))
	structValidator.AddField("Email", user.Email, Required(), MatchRegex("^[a-z0-9._%+-]+@[a-z0-9.-]+\\.[a-z]{2,}$"))
	structValidator.AddField("Age", user.Age, Required())

	errors := structValidator.Validate()
	if errors.HasErrors() {
		t.Errorf("Error in:\n%v\n", errors)
	}
}

func TestConditional(t *testing.T) {
	validateName := func(user User) bool {
		structValidator := &StructValidator{}
		allowAdminName := func() bool {
			return user.Name != "admin"
		}
		structValidator.AddField(
			"Name",
			user.Name,
			Conditional(allowAdminName, MinLength(6)),
		)
		return structValidator.Validate().HasErrors()
	}

	cases := []struct {
		user     User
		hasError bool
	}{
		{User{Name: ""}, true},
		{User{Name: "dalei"}, true},
		{User{Name: "admin"}, false},
		{User{Name: "long name"}, false},
	}
	for _, c := range cases {
		hasError := validateName(c.user)
		if hasError != c.hasError {
			t.Errorf("name `%s` conditional test result %t", c.user.Name, c.hasError)
		}
	}
}

func TestEnums(t *testing.T) {
	validateName := func(user User) bool {
		structValidator := &StructValidator{}
		structValidator.AddField(
			"Name",
			user.Name,
			Enums("admin", "dalei"),
		)
		return structValidator.Validate().HasErrors()
	}

	cases := []struct {
		user     User
		hasError bool
	}{
		{User{Name: ""}, true},
		{User{Name: "dalei"}, false},
		{User{Name: "admin"}, false},
		{User{Name: "long name"}, true},
	}
	for _, c := range cases {
		hasError := validateName(c.user)
		if hasError != c.hasError {
			t.Errorf("name `%s` enums test result %t", c.user.Name, c.hasError)
		}
	}
}
