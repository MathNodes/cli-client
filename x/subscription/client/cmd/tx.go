package cmd

import (
	"bufio"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	hubtypes "github.com/sentinel-official/hub/types"
	"github.com/sentinel-official/hub/x/subscription/types"

	"github.com/sentinel-official/cli-client/context"
	clitypes "github.com/sentinel-official/cli-client/types"
)

func GetTxCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "subscription",
		Short:                      "Subscription related subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		txSubscribeToNode(),
		txSubscribeToPlan(),
		txAddQuota(),
		txUpdateQuota(),
		txCancel(),
	)

	return cmd
}

func txSubscribeToNode() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscribe-to-node [node] [deposit]",
		Short: "Subscribe to a node",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			tc, err := context.NewTxContextFromCmd(cmd)
			if err != nil {
				return err
			}

			reader := bufio.NewReader(cmd.InOrStdin())

			password, fromAddr, err := tc.GetPasswordAndAddress(reader, tc.From)
			if err != nil {
				return err
			}

			nodeAddr, err := hubtypes.NodeAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgSubscribeToNodeRequest(
				fromAddr,
				nodeAddr,
				deposit,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			result, err := tc.SignMessagesAndBroadcastTx(password, msg)
			if err != nil {
				return err
			}

			fmt.Println(result)
			return nil
		},
	}

	clitypes.AddTxFlagsToCmd(cmd)

	return cmd
}

func txSubscribeToPlan() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscribe-to-plan [plan] [denom]",
		Short: "Subscribe to a plan",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			tc, err := context.NewTxContextFromCmd(cmd)
			if err != nil {
				return err
			}

			reader := bufio.NewReader(cmd.InOrStdin())

			password, fromAddr, err := tc.GetPasswordAndAddress(reader, tc.From)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgSubscribeToPlanRequest(
				fromAddr,
				id,
				args[1],
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			result, err := tc.SignMessagesAndBroadcastTx(password, msg)
			if err != nil {
				return err
			}

			fmt.Println(result)
			return nil
		},
	}

	clitypes.AddTxFlagsToCmd(cmd)

	return cmd
}

func txAddQuota() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quota-add [id] [address] [bytes]",
		Short: "Add a quota for subscription",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			tc, err := context.NewTxContextFromCmd(cmd)
			if err != nil {
				return err
			}

			reader := bufio.NewReader(cmd.InOrStdin())

			password, fromAddr, err := tc.GetPasswordAndAddress(reader, tc.From)
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

			bytes, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgAddQuotaRequest(
				fromAddr,
				id,
				accAddr,
				sdk.NewInt(bytes),
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			result, err := tc.SignMessagesAndBroadcastTx(password, msg)
			if err != nil {
				return err
			}

			fmt.Println(result)
			return nil
		},
	}

	clitypes.AddTxFlagsToCmd(cmd)

	return cmd
}

func txUpdateQuota() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quota-update [id] [address] [bytes]",
		Short: "Update a quota for subscription",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			tc, err := context.NewTxContextFromCmd(cmd)
			if err != nil {
				return err
			}

			reader := bufio.NewReader(cmd.InOrStdin())

			password, fromAddr, err := tc.GetPasswordAndAddress(reader, tc.From)
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

			bytes, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateQuotaRequest(
				fromAddr,
				id,
				accAddr,
				sdk.NewInt(bytes),
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			result, err := tc.SignMessagesAndBroadcastTx(password, msg)
			if err != nil {
				return err
			}

			fmt.Println(result)
			return nil
		},
	}

	clitypes.AddTxFlagsToCmd(cmd)

	return cmd
}

func txCancel() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "cancel [id]",
		Short:  "Cancel a subscription",
		Args:   cobra.ExactArgs(1),
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			tc, err := context.NewTxContextFromCmd(cmd)
			if err != nil {
				return err
			}

			reader := bufio.NewReader(cmd.InOrStdin())

			_, fromAddr, err := tc.GetPasswordAndAddress(reader, tc.From)
			if err != nil {
				return err
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgCancelRequest(
				fromAddr,
				id,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return nil
		},
	}

	clitypes.AddTxFlagsToCmd(cmd)

	return cmd
}
