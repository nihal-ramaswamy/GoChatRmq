package dto

type TestConfigDto struct {
	Username     string
	Password     string
	DatabaseName string
}

func NewTestConfigDto(username, password, databaseName string) TestConfigDto {
	return TestConfigDto{
		Username:     username,
		Password:     password,
		DatabaseName: databaseName,
	}
}
