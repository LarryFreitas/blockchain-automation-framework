package pap

import (
	"github.com/PM-Master/policy-machine-go/pip"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/pkg/errors"
	"github.com/usnistgov/blossom/chaincode/model"
	"github.com/usnistgov/blossom/chaincode/ngac/pap/ledger"
	dacpolicy "github.com/usnistgov/blossom/chaincode/ngac/pap/policy/dac"
	rbacpolicy "github.com/usnistgov/blossom/chaincode/ngac/pap/policy/rbac"
	statuspolicy "github.com/usnistgov/blossom/chaincode/ngac/pap/policy/status"
	"time"
)

type LicenseAdmin struct {
	graph pip.Graph
}

func NewLicenseAdmin(ctx contractapi.TransactionContextInterface) (*LicenseAdmin, error) {
	la := &LicenseAdmin{}
	err := la.setup(ctx)
	return la, err
}

func (l *LicenseAdmin) setup(ctx contractapi.TransactionContextInterface) error {
	graph, err := ledger.GetGraph(ctx)
	if err != nil {
		return errors.Wrap(err, "error retrieving ngac graph from ledger")
	}

	l.graph = graph

	return nil
}

func (l *LicenseAdmin) Graph() pip.Graph {
	return l.graph
}

func (l *LicenseAdmin) OnboardLicense(ctx contractapi.TransactionContextInterface, license *model.License) error {
	var (
		licenseOA pip.Node
		err       error
	)

	rbacPolicy := rbacpolicy.NewLicensePolicy(l.graph)
	if licenseOA, err = rbacPolicy.OnboardLicense(license); err != nil {
		return errors.Wrap(err, "error configuring license onboard RBAC policy")
	}

	statusPolicy := statuspolicy.NewLicensePolicy(l.graph)
	if err = statusPolicy.OnboardLicense(licenseOA); err != nil {
		return errors.Wrap(err, "error configuring license onboard Status policy")
	}

	return ledger.UpdateGraphState(ctx, l.graph)
}

func (l *LicenseAdmin) OffboardLicense(ctx contractapi.TransactionContextInterface, licenseID string) error {
	rbacPolicy := rbacpolicy.NewLicensePolicy(l.graph)
	if err := rbacPolicy.OffboardLicense(licenseID); err != nil {
		return errors.Wrap(err, "error configuring license offboard RBAC policy")
	}

	return ledger.UpdateGraphState(ctx, l.graph)
}

func (l *LicenseAdmin) CheckoutLicense(ctx contractapi.TransactionContextInterface, agencyName string, licenseID string,
	keys map[string]time.Time) error {
	dacPolicy := dacpolicy.NewLicensePolicy(l.graph)
	if err := dacPolicy.CheckoutLicense(agencyName, licenseID, keys); err != nil {
		return errors.Wrap(err, "error checking out license under the DAC policy")
	}

	return ledger.UpdateGraphState(ctx, l.graph)
}

func (l *LicenseAdmin) CheckinLicense(ctx contractapi.TransactionContextInterface, agencyName string, licenseID string,
	keys []string) error {
	dacPolicy := dacpolicy.NewLicensePolicy(l.graph)
	if err := dacPolicy.CheckinLicense(agencyName, licenseID, keys); err != nil {
		return errors.Wrap(err, "error checking in license under the DAC policy")
	}

	return ledger.UpdateGraphState(ctx, l.graph)
}
