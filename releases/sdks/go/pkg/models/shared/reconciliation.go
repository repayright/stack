// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"github.com/formancehq/formance-sdk-go/pkg/utils"
	"math/big"
	"time"
)

type Reconciliation struct {
	CreatedAt            time.Time           `json:"createdAt"`
	DriftBalances        map[string]*big.Int `json:"driftBalances"`
	Error                *string             `json:"error,omitempty"`
	ID                   string              `json:"id"`
	LedgerBalances       map[string]*big.Int `json:"ledgerBalances"`
	PaymentsBalances     map[string]*big.Int `json:"paymentsBalances"`
	PolicyID             string              `json:"policyID"`
	ReconciledAtLedger   time.Time           `json:"reconciledAtLedger"`
	ReconciledAtPayments time.Time           `json:"reconciledAtPayments"`
	Status               string              `json:"status"`
}

func (r Reconciliation) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(r, "", false)
}

func (r *Reconciliation) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &r, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *Reconciliation) GetCreatedAt() time.Time {
	if o == nil {
		return time.Time{}
	}
	return o.CreatedAt
}

func (o *Reconciliation) GetDriftBalances() map[string]*big.Int {
	if o == nil {
		return map[string]*big.Int{}
	}
	return o.DriftBalances
}

func (o *Reconciliation) GetError() *string {
	if o == nil {
		return nil
	}
	return o.Error
}

func (o *Reconciliation) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *Reconciliation) GetLedgerBalances() map[string]*big.Int {
	if o == nil {
		return map[string]*big.Int{}
	}
	return o.LedgerBalances
}

func (o *Reconciliation) GetPaymentsBalances() map[string]*big.Int {
	if o == nil {
		return map[string]*big.Int{}
	}
	return o.PaymentsBalances
}

func (o *Reconciliation) GetPolicyID() string {
	if o == nil {
		return ""
	}
	return o.PolicyID
}

func (o *Reconciliation) GetReconciledAtLedger() time.Time {
	if o == nil {
		return time.Time{}
	}
	return o.ReconciledAtLedger
}

func (o *Reconciliation) GetReconciledAtPayments() time.Time {
	if o == nil {
		return time.Time{}
	}
	return o.ReconciledAtPayments
}

func (o *Reconciliation) GetStatus() string {
	if o == nil {
		return ""
	}
	return o.Status
}