package log

import "github.com/jedib0t/go-pretty/v6/text"

type colors [3]text.Colors // 0: prefix, 1: message, 2: extra data

var (
	colorsDebug = colors{
		{text.BgCyan, text.FgHiWhite}, // prefix
		{},                            // message
		{text.FgWhite},                // extra data
	}

	colorsInfo = colors{
		{text.BgBlue, text.FgHiWhite}, // prefix
		{},                            // message
		{text.FgWhite},                // extra data
	}

	colorsSuccess = colors{
		{text.BgGreen, text.FgHiWhite, text.Bold}, // prefix
		{},             // message
		{text.FgWhite}, // extra data
	}

	colorsWarn = colors{
		{text.BgYellow, text.FgHiWhite}, // prefix
		{},                              // message
		{text.FgWhite},                  // extra data
	}

	colorsError = colors{
		{text.BgHiRed, text.FgHiWhite, text.Bold}, // prefix
		{},             // message
		{text.FgWhite}, // extra data
	}

	colorsFatal = colors{
		{text.BgHiRed, text.FgBlack, text.Bold}, // prefix
		{text.BgBlack, text.FgHiRed},            // message
		{text.FgWhite},                          // extra data
	}
)
