package media

import "fmt"

func GetAccountMediaPath(accountID string) string {
	return fmt.Sprintf("uploads/account/%s/", accountID)
}
