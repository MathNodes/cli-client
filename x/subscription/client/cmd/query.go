package cmd

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/sentinel-official/cli-client/context"
	clitypes "github.com/sentinel-official/cli-client/types"
	"github.com/sentinel-official/cli-client/x/subscription/types"
)

var (
	subscriptionHeader = []string{
		"ID",
		"Owner",
		"Plan",
		"Expiry",
		"Denom",
		"Node",
		"Price",
		"Deposit",
		"Free",
		"Status",
	}
	quotaHeader = []string{
		"Address",
		"Allocated",
		"Consumed",
	}
)

func QuerySubscription() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscription [id]",
		Short: "Query a subscription",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			qc, err := context.NewQueryContextFromCmd(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			result, err := qc.QuerySubscription(id)
			if err != nil {
				return err
			}

			var (
				item  = types.NewSubscriptionFromRaw(result)
				table = tablewriter.NewWriter(cmd.OutOrStdout())
			)

			table.SetHeader(subscriptionHeader)
			table.Append(
				[]string{
					fmt.Sprintf("%d", item.ID),
					item.Owner,
					fmt.Sprintf("%d", item.Plan),
					item.Expiry.String(),
					item.Denom,
					item.Node,
					item.Price.Raw().String(),
					item.Deposit.Raw().String(),
					clitypes.ToReadableBytes(item.Free, 2),
					item.Status,
				},
			)

			table.Render()
			return nil
		},
	}

	clitypes.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QuerySubscriptions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscriptions",
		Short: "Query subscriptions",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			qc, err := context.NewQueryContextFromCmd(cmd)
			if err != nil {
				return err
			}

			accAddr, err := clitypes.GetAccAddressFromCmd(cmd)
			if err != nil {
				return err
			}

			status, err := clitypes.GetStatusFromCmd(cmd)
			if err != nil {
				return err
			}

			pagination, err := clitypes.GetPageRequestFromCmd(cmd)
			if err != nil {
				return err
			}

			var items types.Subscriptions
			if accAddr != nil {
				result, err := qc.QuerySubscriptionsForAddress(
					accAddr,
					status,
					pagination,
				)
				if err != nil {
					return err
				}

				items = append(items, types.NewSubscriptionsFromRaw(result)...)
			} else {
				result, err := qc.QuerySubscriptions(pagination)
				if err != nil {
					return err
				}

				items = append(items, types.NewSubscriptionsFromRaw(result)...)
			}

			table := tablewriter.NewWriter(cmd.OutOrStdout())
			table.SetHeader(subscriptionHeader)

			for i := 0; i < len(items); i++ {
				table.Append(
					[]string{
						fmt.Sprintf("%d", items[i].ID),
						items[i].Owner,
						fmt.Sprintf("%d", items[i].Plan),
						items[i].Expiry.String(),
						items[i].Denom,
						items[i].Node,
						items[i].Price.Raw().String(),
						items[i].Deposit.Raw().String(),
						clitypes.ToReadableBytes(items[i].Free, 2),
						items[i].Status,
					},
				)
			}

			table.Render()
			return nil
		},
	}

	clitypes.AddQueryFlagsToCmd(cmd)
	clitypes.AddPaginationFlagsToCmd(cmd, "subscriptions")

	cmd.Flags().String(clitypes.FlagAddress, "", "filter with account address")
	cmd.Flags().String(clitypes.FlagStatus, "Active", "filter with status (Active|Inactive)")

	return cmd
}

func QueryQuota() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quota [subscription] [address]",
		Short: "Query a quota",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			qc, err := context.NewQueryContextFromCmd(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			accAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			result, err := qc.QueryQuota(
				id,
				accAddr,
			)
			if err != nil {
				return err
			}

			var (
				item  = types.NewQuotaFromRaw(result)
				table = tablewriter.NewWriter(cmd.OutOrStdout())
			)

			table.SetHeader(quotaHeader)
			table.Append(
				[]string{
					item.Address,
					clitypes.ToReadableBytes(item.Allocated, 2),
					clitypes.ToReadableBytes(item.Consumed, 2),
				},
			)

			table.Render()
			return nil
		},
	}

	clitypes.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QueryQuotas() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quotas [subscription]",
		Short: "Query quotas of a subscription",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			qc, err := context.NewQueryContextFromCmd(cmd)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			pagination, err := clitypes.GetPageRequestFromCmd(cmd)
			if err != nil {
				return err
			}

			result, err := qc.QueryQuotas(
				id,
				pagination,
			)
			if err != nil {
				return err
			}

			var (
				items = types.NewQuotasFromRaw(result)
				table = tablewriter.NewWriter(cmd.OutOrStdout())
			)

			table.SetHeader(quotaHeader)
			for i := 0; i < len(items); i++ {
				table.Append(
					[]string{
						items[i].Address,
						clitypes.ToReadableBytes(items[i].Allocated, 2),
						clitypes.ToReadableBytes(items[i].Consumed, 2),
					},
				)
			}

			table.Render()
			return nil
		},
	}

	clitypes.AddQueryFlagsToCmd(cmd)
	clitypes.AddPaginationFlagsToCmd(cmd, "quotas")

	return cmd
}
