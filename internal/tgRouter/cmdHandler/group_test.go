package cmdhandler

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ponyCorp/rebornPony/internal/repository"
	"github.com/ponyCorp/rebornPony/internal/services/sender"
)

func TestCmdHandler_NewGroup(t *testing.T) {
	type fields struct {
		mapRouter map[string]rout
		rep       *repository.Repository
		rules     *rules
		sender    *sender.Sender
	}
	type args struct {
		groupName string
	}

	t.Run("TestCmdHandler_NewGroup", func(t *testing.T) {
		h := NewCmdHandler(nil, nil)
		testGroup := h.NewGroup("test")
		testGroup.Handle("test", "test", func(update *tgbotapi.Update, cmd, arg string) {

		})

		// if got := h.NewGroup(tt.args.groupName); !reflect.DeepEqual(got, tt.want) {
		// 	t.Errorf("CmdHandler.NewGroup() = %v, want %v", got, tt.want)
		// }
	})

}
