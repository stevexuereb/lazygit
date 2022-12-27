package custom_commands

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// NOTE: we're getting a weird offset in the popup prompt for some reason. Not sure what's behind that.

var MenuFromCommand = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Using menuFromCommand prompt type",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupRepo: func(shell *Shell) {
		shell.
			EmptyCommit("foo").
			EmptyCommit("bar").
			EmptyCommit("baz").
			NewBranch("feature/foo")
	},
	SetupConfig: func(cfg *config.AppConfig) {
		cfg.UserConfig.CustomCommands = []config.CustomCommand{
			{
				Key:     "a",
				Context: "localBranches",
				Command: `echo "{{index .PromptResponses 0}} {{index .PromptResponses 1}} {{ .SelectedLocalBranch.Name }}" > output.txt`,
				Prompts: []config.CustomCommandPrompt{
					{
						Type:        "menuFromCommand",
						Title:       "Choose commit message",
						Command:     `git log --oneline --pretty=%B`,
						Filter:      `(?P<commit_message>.*)`,
						ValueFormat: `{{ .commit_message }}`,
						LabelFormat: `{{ .commit_message | yellow }}`,
					},
					{
						Type:         "input",
						Title:        "Description",
						InitialValue: `{{ if .SelectedLocalBranch.Name }}Branch: #{{ .SelectedLocalBranch.Name }}{{end}}`,
					},
				},
			},
		}
	},
	Run: func(
		shell *Shell,
		t *TestDriver,
		keys config.KeybindingConfig,
	) {
		t.Model().WorkingTreeFileCount(0)
		t.Views().Branches().
			Focus().
			Press("a")

		t.ExpectMenu().Title(Equals("Choose commit message")).Select(Contains("bar")).Confirm()

		t.ExpectPrompt().Title(Equals("Description")).Type(" my branch").Confirm()

		t.Model().WorkingTreeFileCount(1)

		t.Views().Files().Focus().SelectedLine(Contains("output.txt"))
		t.Views().Main().Content(Contains("bar Branch: #feature/foo my branch feature/foo"))
	},
})
