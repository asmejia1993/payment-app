package util

import m "github.com/asmejia1993/payment-app/db/sqlc"

func IsSupportedUserType(uType string) bool {
	switch uType {
	case m.UserTypeEnum.Customer, m.UserTypeEnum.Merchant:
		return true
	}
	return false
}
