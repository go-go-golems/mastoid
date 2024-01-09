package cmds

import (
	"context"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/layers"
	"github.com/go-go-golems/glazed/pkg/cmds/parameters"
	"github.com/go-go-golems/glazed/pkg/middlewares"
	"github.com/go-go-golems/glazed/pkg/settings"
	"github.com/go-go-golems/glazed/pkg/types"
	"github.com/go-go-golems/mastoid/cmd/mastoid/pkg"
	"github.com/mattn/go-mastodon"
	"github.com/pkg/errors"
)

type ThreadCmd struct {
	*cmds.CommandDescription
}

func NewThreadCommand() (*ThreadCmd, error) {
	glazedParameterLayer, err := settings.NewGlazedParameterLayers()
	if err != nil {
		return nil, errors.Wrap(err, "could not create Glazed parameter layer")
	}

	return &ThreadCmd{
		CommandDescription: cmds.NewCommandDescription(
			"thread",
			cmds.WithShort("Output thread as structured data"),
			cmds.WithFlags(
				parameters.NewParameterDefinition(
					"status",
					parameters.ParameterTypeString,
					parameters.WithHelp("Thread status id"),
					parameters.WithShortFlag("s"),
				),
				parameters.NewParameterDefinition(
					"verbose",
					parameters.ParameterTypeBool,
					parameters.WithHelp("Verbose output"),
					parameters.WithShortFlag("v"),
					parameters.WithDefault(false),
				),
			),
			cmds.WithLayersList(glazedParameterLayer),
		),
	}, nil
}

type ThreadSettings struct {
	Status  string `glazed.parameter:"status"`
	Verbose bool   `glazed.parameter:"verbose"`
}

var _ cmds.GlazeCommand = (*ThreadCmd)(nil)

func (c *ThreadCmd) RunIntoGlazeProcessor(
	ctx context.Context,
	parsedLayers *layers.ParsedLayers,
	gp middlewares.Processor,
) error {
	s := &ThreadSettings{}

	err := parsedLayers.InitializeStruct(layers.DefaultSlug, s)
	if err != nil {
		return err
	}

	statusID := pkg.ExtractID(s.Status)
	if statusID == "" {
		return errors.New("no status ID provided")
	}

	credentials, err := pkg.LoadCredentials()
	if err != nil {
		return errors.Wrap(err, "could not load credentials")
	}

	client, err := pkg.CreateClientAndAuthenticate(ctx, credentials)
	if err != nil {
		return errors.Wrap(err, "could not create client")
	}

	status_, err := client.GetStatus(ctx, mastodon.ID(statusID))
	if err != nil {
		return err
	}

	context_, err := client.GetStatusContext(ctx, status_.ID)
	if err != nil {
		return err
	}

	thread := &pkg.Thread{
		Nodes: map[mastodon.ID]*pkg.Node{},
	}

	thread.AddStatus(status_)
	thread.AddContextAndGetMissingIDs(status_.ID, context_)

	printNode := func(node *pkg.Node, depth int, siblingIdx int) error {
		var row types.Row
		if s.Verbose {
			row = types.NewRowFromStruct(node.Status, true)
		} else {
			row = types.NewRow(
				types.MRP("ID", node.Status.ID),
				types.MRP("CreatedAt", node.Status.CreatedAt),
				types.MRP("Author", node.Status.Account.Acct),
				types.MRP("Content", node.Status.Content),
				types.MRP("InReplyToID", node.Status.InReplyToID),
			)
		}
		row.Set("Depth", depth)
		row.Set("SiblingIndex", siblingIdx)
		err = gp.AddRow(ctx, row)
		if err != nil {
			return err
		}
		return nil
	}

	err = thread.WalkDepthFirst(printNode)
	if err != nil {
		return err
	}

	return nil
}
