package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

// taken from here: https://github.com/spf13/cobra/blob/master/shell_completions.md
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

$ source <(deploy1 completion bash)

# To load completions for each session, execute once:
Linux:
  $ deploy1 completion bash > /etc/bash_completion.d/deploy1
MacOS:
  $ deploy1 completion bash > /usr/local/etc/bash_completion.d/deploy1

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ deploy1 completion zsh > "${fpath[1]}/_deploy1"

# You will need to start a new shell for this setup to take effect.

Fish:

$ deploy1 completion fish | source

# To load completions for each session, execute once:
$ deploy1 completion fish > ~/.config/fish/completions/deploy1.fish
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			return cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			return cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			return cmd.Root().GenPowerShellCompletion(os.Stdout)
		default:
			return nil
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
