package utils

import (
	"github.com/chelnak/ysmrr"
	"github.com/chelnak/ysmrr/pkg/animations"
	"github.com/chelnak/ysmrr/pkg/colors"
)

var Spinner = ysmrr.NewSpinnerManager(
	ysmrr.WithAnimation(animations.Dots),
	ysmrr.WithSpinnerColor(colors.FgHiBlue),
)
