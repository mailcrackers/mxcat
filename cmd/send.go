package cmd

import (
	"mxcat/internal/smtp"

	"github.com/spf13/cobra"
)

var options = smtp.Options{}
var send = cobra.Command{
	Use:   "send",
	Short: "This command is used to send emails from CLI",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	send.Flags().StringVar(&options.Host, "host", "127.0.0.1", "target server hostname or IP address")
	send.Flags().IntVar(&options.Port, "port", 25, "port to connect to on the target server")
	send.Flags().BoolVar(&options.UseTLS, "tls", false, "use TLS connection")
	send.Flags().BoolVar(&options.StartTLS, "start-tls", false, "start session with STARTTLS command")
	send.Flags().StringVar(&options.From, "from", "", "sender email address used in the MAIL FROM command")
	send.Flags().StringVar(&options.Helo, "helo", "", "domain name to use in HELO command (used to identify this client)")
	send.Flags().StringVar(&options.Ehlo, "ehlo", "", "domain name to use in EHLO command (used to identify this client)")
	send.Flags().StringVar(&options.Lhlo, "lhlo", "", "domain name to use in LHLO command (used to identify this client)")
	send.Flags().StringSliceVar(&options.To, "to", []string{}, "recipient email addresses used in RCPT TO commands")

	cmd.AddCommand(&send)
}
