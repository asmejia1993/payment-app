package bank

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/plaid/plaid-go/v20/plaid"
)

type plaidConfig struct {
	plaidClient *plaid.APIClient
}

func NewPlaidClient(plaidClientId, plaidSecret string) Simulator {
	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", plaidClientId)
	configuration.AddDefaultHeader("PLAID-SECRET", plaidSecret)
	configuration.UseEnvironment(plaid.Sandbox)
	client := plaid.NewAPIClient(configuration)

	return &plaidConfig{plaidClient: client}
}

// RefundPayment implements Simulator.
func (p *plaidConfig) RefundPayment(ctx context.Context) {
	id := uuid.New().String()
	fmt.Printf("uuid: %s\n", id)
	request := plaid.NewSandboxTransferRefundSimulateRequest("TRANSFER_ID", "posted")
	transferSimulateResp, _, err := p.plaidClient.PlaidApi.SandboxTransferRefundSimulate(ctx).SandboxTransferRefundSimulateRequest(
		*request,
	).Execute()
	if err != nil {
		fmt.Printf("error on simulate transfer: %v", err)
	}
	fmt.Println(transferSimulateResp)
}

// StartPaymentProcess implements Simulator.
func (p *plaidConfig) StartPaymentProcess(ctx context.Context) (string, error) {
	request := plaid.NewSandboxTransferSimulateRequest("TRANSFER_ID", "posted")
	transferSimulateResp, _, err := p.plaidClient.PlaidApi.SandboxTransferSimulate(ctx).SandboxTransferSimulateRequest(
		*request,
	).Execute()
	if err != nil {
		fmt.Printf("error on simulate transfer: %v", err)
	}

	fmt.Println(transferSimulateResp)
	return transferSimulateResp.RequestId, nil
}
