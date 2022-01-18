package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	var err error

	result := make(DomainStat)

	scanner := bufio.NewScanner(r)

	subDomain := "." + strings.ToLower(domain)

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.Contains(line, subDomain) {
			continue
		}

		var user User
		if err = json.Unmarshal([]byte(line), &user); err != nil {
			return nil, err
		}

		userEmailDomain := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])

		if !strings.HasSuffix(userEmailDomain, subDomain) {
			continue
		}

		result[userEmailDomain]++
	}

	return result, nil
}
