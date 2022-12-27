package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var Staged = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Staging a couple files, going in the staged files menu, unstaging a line then committing",
	ExtraCmdArgs: "",
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.
			CreateFile("myfile", "myfile content\nwith a second line").
			CreateFile("myfile2", "myfile2 content")
	},
	Run: func(shell *Shell, t *TestDriver, keys config.KeybindingConfig) {
		t.Model().CommitCount(0)

		t.Views().Files().
			IsFocused().
			SelectedLine(Contains("myfile")).
			PressPrimaryAction(). // stage the file
			PressEnter()

		t.Views().StagingSecondary().
			IsFocused().
			Tap(func() {
				// we start with both lines having been staged
				t.Views().StagingSecondary().Content(Contains("+myfile content"))
				t.Views().StagingSecondary().Content(Contains("+with a second line"))
				t.Views().Staging().Content(DoesNotContain("+myfile content"))
				t.Views().Staging().Content(DoesNotContain("+with a second line"))
			}).
			// unstage the selected line
			PressPrimaryAction().
			Tap(func() {
				// the line should have been moved to the main view
				t.Views().StagingSecondary().Content(DoesNotContain("+myfile content"))
				t.Views().StagingSecondary().Content(Contains("+with a second line"))
				t.Views().Staging().Content(Contains("+myfile content"))
				t.Views().Staging().Content(DoesNotContain("+with a second line"))
			}).
			Press(keys.Files.CommitChanges)

		commitMessage := "my commit message"
		t.ExpectCommitMessagePanel().Type(commitMessage).Confirm()

		t.Model().CommitCount(1)
		t.Model().HeadCommitMessage(Equals(commitMessage))
		t.Views().StagingSecondary().IsFocused()

		// TODO: assert that the staging panel has been refreshed (it currently does not get correctly refreshed)
	},
})
