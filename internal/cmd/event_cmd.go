package cmd

import (
	"fmt"

	"github.com/objectisundefined/github-cli/internal/client"
	"github.com/spf13/cobra"
)

var (
	eventsLimit     int
	eventsSinceTime string
	eventsOffsetDur string
)

func NewEventsCommand() *cobra.Command {
	m := &cobra.Command{
		Use:   "events [users]",
		Short: "Github CLI for tracing user events",
		Args:  cobra.MinimumNArgs(0),
		Run:   runEventsCommandFunc,
	}

	m.Flags().StringVar(&eventsSinceTime, "since", "", fmt.Sprintf("Issue Since Time, format is %s", TimeFormat))
	m.Flags().StringVar(&eventsOffsetDur, "offset", "-336h", "The offset of since time")
	m.Flags().IntVar(&eventsLimit, "limit", 20, "Maximum issues limit for a repository")

	return m
}

func runEventsCommandFunc(cmd *cobra.Command, args []string) {
	user := cfg.Account
	if len(args) > 0 {
		user = args[0]
	}

	opts := NewEventOptions()
	opts.Limit = eventsLimit

	opts.RangeTime.adjust(eventsSinceTime, eventsOffsetDur)

	ec := EventClient{client.NewClient(ctx, cfg)}
	events, err := ec.ListEventsByUser(ctx, user, opts)

	if err != nil {
		panic(err)
	}

	for _, event := range events {
		fmt.Printf("%s - %s\n", event.GetCreatedAt().Format(TimeFormat), formatEvent(event))
	}
}
