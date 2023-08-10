package suite

import (
	"math/big"
	"time"

	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Given("some empty environment", func() {
	When("creating a transaction on a ledger", func() {
		var (
			timestamp = time.Now().Round(time.Second).UTC()
			rsp       *shared.TransactionsResponse
		)
		BeforeEach(func() {
			// Create a transaction
			response, err := Client().Transactions.CreateTransaction(
				TestContext(),
				operations.CreateTransactionRequest{
					PostTransaction: shared.PostTransaction{
						Metadata: map[string]any{},
						Postings: []shared.Posting{
							{
								Amount:      big.NewInt(100),
								Asset:       "USD",
								Source:      "world",
								Destination: "alice",
							},
						},
						Timestamp: &timestamp,
					},
					Ledger: "default",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(200))

			rsp = response.TransactionsResponse

			// Check existence on api
			getResponse, err := Client().Transactions.GetTransaction(
				TestContext(),
				operations.GetTransactionRequest{
					Ledger: "default",
					Txid:   rsp.Data[0].Txid,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(getResponse.StatusCode).To(Equal(200))
		})
		It("should fail if the transaction does not exist", func() {
			metadata := map[string]any{
				"foo": "bar",
			}

			response, err := Client().Transactions.AddMetadataOnTransaction(
				TestContext(),
				operations.AddMetadataOnTransactionRequest{
					RequestBody: metadata,
					Ledger:      "default",
					Txid:        666,
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(404))
		})
		Then("adding a metadata", func() {
			metadata := map[string]any{
				"foo": "bar",
			}
			BeforeEach(func() {
				response, err := Client().Transactions.AddMetadataOnTransaction(
					TestContext(),
					operations.AddMetadataOnTransactionRequest{
						RequestBody: metadata,
						Ledger:      "default",
						Txid:        rsp.Data[0].Txid,
					},
				)
				Expect(err).To(Succeed())
				Expect(response.StatusCode).To(Equal(204))
			})
			It("should eventually be available on api", func() {
				// Check existence on api
				response, err := Client().Transactions.GetTransaction(
					TestContext(),
					operations.GetTransactionRequest{
						Ledger: "default",
						Txid:   rsp.Data[0].Txid,
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(200))

				Expect(response.TransactionResponse.Data.Metadata).Should(Equal(metadata))
			})
		})
	})
})
