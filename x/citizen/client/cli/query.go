package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"uagd/x/citizen/types"
)

func CmdRegionByAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "region [address]",
		Short: "Query region by address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.RegionByAddress(context.Background(), &types.QueryRegionByAddressRequest{Address: args[0]})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdAddressesByRegion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addresses-by-region [region-id]",
		Short: "List addresses by region",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.AddressesByRegion(context.Background(), &types.QueryAddressesByRegionRequest{RegionId: args[0]})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Show the citizen module parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
