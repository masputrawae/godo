package password

import "golang.org/x/crypto/bcrypt"

func Hash(p string) (*string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return new(string(b)), nil
}

func Check(h, p string) bool {
	return bcrypt.CompareHashAndPassword([]byte(h), []byte(p)) == nil
}
