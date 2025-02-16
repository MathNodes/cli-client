package cmd

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/go-kit/kit/transport/http/jsonrpc"
	hubtypes "github.com/sentinel-official/hub/types"
	sessiontypes "github.com/sentinel-official/hub/x/session/types"
	"github.com/spf13/cobra"

	"github.com/sentinel-official/cli-client/context"
	wireguardtypes "github.com/sentinel-official/cli-client/services/wireguard/types"
	clitypes "github.com/sentinel-official/cli-client/types"
)

func parseResolversFromCmd(cmd *cobra.Command) ([]net.IP, error) {
	v, err := cmd.Flags().GetStringArray(clitypes.FlagResolver)
	if err != nil {
		return nil, err
	}

	items := make([]net.IP, 0, len(v))
	for _, s := range v {
		item := net.ParseIP(s)
		if item == nil {
			return nil, fmt.Errorf("invalid resolver ip %s", s)
		}

		items = append(items, item)
	}

	return items, nil
}

func ConnectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connect [subscription] [address]",
		Short: "Connect to a node",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			nodeAddr, err := hubtypes.NodeAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			rating, err := cmd.Flags().GetUint64(clitypes.FlagRating)
			if err != nil {
				return err
			}

			resolvers, err := parseResolversFromCmd(cmd)
			if err != nil {
				return err
			}

			tc, err := context.NewTxContextFromCmd(cmd)
			if err != nil {
				return err
			}

			sc, err := context.NewServiceContextFromCmd(cmd)
			if err != nil {
				return err
			}

			status, err := sc.GetStatus()
			if err != nil {
				return err
			}

			if status.IFace != "" {
				if err := sc.Disconnect(); err != nil {
					return err
				}
			}

			var (
				messages []sdk.Msg
				reader   = bufio.NewReader(cmd.InOrStdin())
			)

			password, from, err := tc.GetPasswordAndAddress(reader, tc.From)
			if err != nil {
				return err
			}

			session, err := tc.QueryActiveSession(from)
			if err != nil {
				return err
			}

			if session != nil {
				messages = append(
					messages,
					sessiontypes.NewMsgEndRequest(
						from,
						session.Id,
						rating,
					),
				)
			}

			messages = append(
				messages,
				sessiontypes.NewMsgStartRequest(
					from,
					id,
					nodeAddr,
				),
			)

			txRes, err := tc.SignMessagesAndBroadcastTx(password, messages...)
			if err != nil {
				return err
			}

			fmt.Println(txRes)

			session, err = tc.QueryActiveSession(from)
			if err != nil {
				return err
			}
			if session == nil {
				return fmt.Errorf("active session does not exist for subscription %d", id)
			}

			node, err := tc.QueryNode(nodeAddr)
			if err != nil {
				return err
			}

			wgPrivateKey, err := wireguardtypes.NewPrivateKey()
			if err != nil {
				return err
			}

			signMsgRes, err := tc.SignMessage(
				password,
				tc.From,
				sdk.Uint64ToBigEndian(session.Id),
			)
			if err != nil {
				return err
			}

			signature, err := base64.StdEncoding.DecodeString(signMsgRes.Signature)
			if err != nil {
				return err
			}

			var (
				resp     clitypes.RestResponseBody
				endpoint = fmt.Sprintf(
					"%s/accounts/%s/sessions/%d",
					strings.Trim(node.RemoteURL, "/"), from, session.Id,
				)
			)

			buf, err := json.Marshal(
				map[string]interface{}{
					"key":       wgPrivateKey.Public().String(),
					"signature": signature,
				},
			)
			if err != nil {
				return err
			}

			r, err := tc.Post(endpoint, jsonrpc.ContentType, bytes.NewBuffer(buf))
			if err != nil {
				return err
			}

			if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
				return err
			}
			if resp.Error != nil {
				return fmt.Errorf(resp.Error.Message)
			}

			info, err := base64.StdEncoding.DecodeString(resp.Result.(string))
			if err != nil {
				return err
			}

			return sc.Connect(
				tc.Backend,
				password,
				from.String(),
				nodeAddr.String(),
				session.Id,
				info,
				[][]byte{wgPrivateKey.Bytes()},
				resolvers,
			)
		},
	}

	clitypes.AddServiceFlagsToCmd(cmd)
	clitypes.AddTxFlagsToCmd(cmd)

	cmd.Flags().StringArray(clitypes.FlagResolver, nil, "provide additional DNS servers")
	cmd.Flags().Uint64(clitypes.FlagRating, 0, "rate the session quality between 0 and 10")

	return cmd
}
