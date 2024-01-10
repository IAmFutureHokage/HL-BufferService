package toModelTelegramConverter

import (
	"fmt"

	"github.com/IAmFutureHokage/HL-BufferService/internal/app/model"
	types "github.com/IAmFutureHokage/HL-BufferService/pkg/types"
)

func toModel(mainTg *model.Telegram, draftTg *types.Telegram) error {

	if mainTg == nil || draftTg == nil {
		return fmt.Errorf("data is not valid")
	}

	return nil
}
