package account

type Accounts struct {
	Accounts []*Account
}

func (a *Accounts) AddAccount(account *Account) {
	a.Accounts = append(a.Accounts, account)
}
