package organization

import (
	"fmt"
	"ocm-support-cli/pkg/capability"
	"ocm-support-cli/pkg/label"
	"ocm-support-cli/pkg/quota"

	sdk "github.com/openshift-online/ocm-sdk-go"
	v1 "github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1"

	"ocm-support-cli/pkg/subscription"
	"ocm-support-cli/pkg/types"
)

type Organization struct {
	types.Meta
	Name          string
	Subscriptions []subscription.Subscription `json:",omitempty"`
	Quota         []quota.Quota               `json:",omitempty"`
	Labels        label.LabelsList            `json:",omitempty"`
	Capabilities  capability.CapabilityList   `json:",omitempty"`
}

func GetOrganizations(key string, limit int, fetchLabels bool, fetchCapabilities bool, conn *sdk.Connection) ([]*v1.Organization, error) {
	search := fmt.Sprintf("id = '%s'", key)
	search += fmt.Sprintf("or external_id = '%s'", key)
	search += fmt.Sprintf("or ebs_account_id = '%s'", key)

	organizations, err := conn.AccountsMgmt().V1().Organizations().List().Parameter("fetchLabels", fetchLabels).Parameter("fetchCapabilities", fetchCapabilities).Size(limit).Search(search).Send()
	if err != nil {
		return []*v1.Organization{}, fmt.Errorf("can't retrieve organizations: %w", err)
	}

	return organizations.Items().Slice(), nil
}

func PresentOrganization(organization *v1.Organization, subscriptions []*v1.Subscription, quotaCostList []*v1.QuotaCost) Organization {
	return Organization{
		Meta:          types.Meta{ID: organization.ID(), HREF: organization.HREF()},
		Name:          organization.Name(),
		Subscriptions: subscription.PresentSubscriptions(subscriptions),
		Quota:         quota.PresentQuotaList(quotaCostList),
		Labels:        label.PresentLabels(organization.Labels()),
		Capabilities:  capability.PresentCapabilities(organization.Capabilities()),
	}
}