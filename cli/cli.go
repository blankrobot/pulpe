package cli

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/blankrobot/pulpe"
	"github.com/blankrobot/pulpe/http"
	"github.com/blankrobot/pulpe/mock"
	"github.com/spf13/cobra"
)

// New returns the pulpe CLI application.
func New() *cobra.Command {
	cmd := cobra.Command{
		Use: "pulpe",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(NewServerCmd())
	cmd.AddCommand(NewVersionCmd())
	return &cmd
}

// NewVersionCmd returns a command that displays the pulpe version number.
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:  "version",
		Long: "Display the version number",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(pulpe.Version)
			return nil
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}
}

// NewServerCmd returns a ServerCmd.
func NewServerCmd() *cobra.Command {
	var s ServerCmd

	cmd := cobra.Command{
		Use:           "server",
		RunE:          s.Run,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Flags().StringVar(&s.addr, "http", ":4000", "HTTP address")

	return &cmd
}

// ServerCmd is a command the runs the pulpe server.
type ServerCmd struct {
	addr string
}

// Run creates a bolt client and runs the HTTP server.
func (c *ServerCmd) Run(cmd *cobra.Command, args []string) error {
	client := mock.NewClient()
	srv := http.NewServer(c.addr, http.NewHandler(client))

	err := srv.Open()
	if err != nil {
		return err
	}

	log.Printf("Serving HTTP on address %s\n", c.addr)

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	<-ch
	fmt.Println()
	log.Println("Stopping server...")
	err = srv.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("OK")
	return nil
}