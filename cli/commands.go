package cli

import "github.com/codegangsta/cli"

var (
	commands = []cli.Command{
		{
			Name:      "create",
			ShortName: "c",
			Usage:     "Create a cluster",
			Action:    create,
		},
		{
			Name:      "list",
			ShortName: "l",
			Usage:     "List nodes in a cluster",
			Flags:     []cli.Flag{flTimeout, flDiscoveryOpt},
			Action:    list,
		},
		{
			Name:      "manage",
			ShortName: "m",
			Usage:     "Manage a docker cluster",
			Flags: []cli.Flag{
				flStrategy, flFilter,
				flHosts,
				flLeaderElection, flLeaderTTL, flManageAdvertise,
				flTLS, flTLSCaCert, flTLSCert, flTLSKey, flTLSVerify,
				flHeartBeat,
				flEnableCors,
				flCluster, flDiscoveryOpt, flClusterOpt},
			Action: manage,
		},
		{
			Name:      "join",
			ShortName: "j",
			Usage:     "join a docker cluster",
			Flags:     []cli.Flag{flJoinAdvertise, flHeartBeat, flTTL, flDiscoveryOpt},
			Action:    join,
		},
		{
			Name:      "proxy",
			ShortName: "p",
			Usage:     "proxy to a docker ",
			Flags: []cli.Flag{
				flHosts,
				flTLS, flTLSCaCert, flTLSCert, flTLSKey, flTLSVerify,
				flCluster, flClusterOpt, flStrategy, flFilter,
				flJoinAdvertise, flHeartBeat, flTTL,
			},
			Action: proxy,
		},
	}
)
