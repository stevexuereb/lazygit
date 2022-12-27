package misc

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var ConfirmOnQuit = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Quitting with a confirm prompt",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig: func(config *config.AppConfig) {
		config.UserConfig.ConfirmOnQuit = true
	},
	SetupRepo: func(shell *Shell) {},
	Run: func(shell *Shell, t *TestDriver, keys config.KeybindingConfig) {
		t.Model().CommitCount(0)

		t.Views().Files().
			IsFocused().
			Press(keys.Universal.Quit)

		t.ExpectConfirmation().
			Title(Equals("")).
			Content(Contains("Are you sure you want to quit?")).
			Confirm()
	},
})
