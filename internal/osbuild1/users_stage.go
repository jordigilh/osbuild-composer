package osbuild1

import (
	"github.com/osbuild/osbuild-composer/internal/blueprint"
	"github.com/osbuild/osbuild-composer/internal/crypt"
)

type UsersStageOptions struct {
	Users map[string]UsersStageOptionsUser `json:"users"`
}

func (UsersStageOptions) isStageOptions() {}

type UsersStageOptionsUser struct {
	UID         *int     `json:"uid,omitempty"`
	GID         *int     `json:"gid,omitempty"`
	Groups      []string `json:"groups,omitempty"`
	Description *string  `json:"description,omitempty"`
	Home        *string  `json:"home,omitempty"`
	Shell       *string  `json:"shell,omitempty"`
	Password    *string  `json:"password,omitempty"`
	Key         *string  `json:"key,omitempty"`
}

func NewUsersStage(options *UsersStageOptions) *Stage {
	return &Stage{
		Name:    "org.osbuild.users",
		Options: options,
	}
}

func NewUsersStageOptions(userCustomizations []blueprint.UserCustomization) (*UsersStageOptions, error) {
	if len(userCustomizations) == 0 {
		return nil, nil
	}

	users := make(map[string]UsersStageOptionsUser, len(userCustomizations))
	for _, uc := range userCustomizations {
		if uc.Password != nil && !crypt.PasswordIsCrypted(*uc.Password) {
			cryptedPassword, err := crypt.CryptSHA512(*uc.Password)
			if err != nil {
				return nil, err
			}

			uc.Password = &cryptedPassword
		}

		user := UsersStageOptionsUser{
			UID:         uc.UID,
			GID:         uc.GID,
			Groups:      uc.Groups,
			Description: uc.Description,
			Home:        uc.Home,
			Shell:       uc.Shell,
			Password:    uc.Password,
			Key:         uc.Key,
		}
		users[uc.Name] = user
	}

	return &UsersStageOptions{Users: users}, nil
}
