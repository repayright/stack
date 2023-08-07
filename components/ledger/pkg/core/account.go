package core

import (
	"encoding/json"
	"regexp"

	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/uptrace/bun"
)

const (
	WORLD = "world"
)

type Account struct {
	bun.BaseModel `bun:"table:accounts,alias:accounts"`

	Address  string            `json:"address"`
	Metadata metadata.Metadata `json:"metadata"`
}

func (a Account) copy() Account {
	a.Metadata = a.Metadata.Copy()
	return a
}

func NewAccount(address string) Account {
	return Account{
		Address:  address,
		Metadata: metadata.Metadata{},
	}
}

type AccountWithVolumes struct {
	Account `bun:",extend"`
	Volumes VolumesByAssets `json:"volumes,omitempty" bun:"volumes,type:jsonb"`
	EffectiveVolumes VolumesByAssets `json:"effectiveVolumes" bun:"effectiveVolumes,type:jsonb"`
}

func NewAccountWithVolumes(address string) *AccountWithVolumes {
	return &AccountWithVolumes{
		Account: Account{
			Address:  address,
			Metadata: metadata.Metadata{},
		},
		Volumes: map[string]*Volumes{},
	}
}

func (v AccountWithVolumes) MarshalJSON() ([]byte, error) {
	type aux AccountWithVolumes
	return json.Marshal(struct {
		aux
		Balances BalancesByAssets `json:"balances"`
	}{
		aux:      aux(v),
		Balances: v.Volumes.Balances(),
	})
}

func (v AccountWithVolumes) Copy() AccountWithVolumes {
	v.Account = v.Account.copy()
	v.Volumes = v.Volumes.copy()
	return v
}

const AccountPattern = "^[a-zA-Z_]+[a-zA-Z0-9_:]*$"

var AccountRegexp = regexp.MustCompile(AccountPattern)
