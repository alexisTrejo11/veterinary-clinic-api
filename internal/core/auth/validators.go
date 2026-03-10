package auth

func ValidatePasswordStrength(plainPassword string) error {
	if len(plainPassword) < 8 {
		return nil
	}
	return nil
}
