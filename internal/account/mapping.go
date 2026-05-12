package account

// takes an account and returns the account dto
func toAccountDTO(acc *account) accountDTO {
	return accountDTO{
		ID:                  acc.ID,
		Username:            acc.Username,
		Active:              acc.Active,
		FailedLoginAttempts: acc.FailedLoginAttempts,
		LockedUntil:         acc.LockedUntil,
		CreatedAt:           acc.CreatedAt,
	}
}
