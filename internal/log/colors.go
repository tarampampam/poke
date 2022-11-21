package log

import "github.com/jedib0t/go-pretty/v6/text"

type colors [3]text.Colors // 0: prefix, 1: message, 2: extra data

var (
	colorsDebug = colors{
		{text.BgWhite}, // prefix
		{},             // message
		{text.FgWhite}, // extra data
	}

	colorsInfo = colors{
		{text.BgHiBlue, text.FgHiWhite}, // prefix
		{},                              // message
		{text.FgWhite},                  // extra data
	}

	colorsSuccess = colors{
		{text.BgHiGreen, text.FgHiWhite}, // prefix
		{},                               // message
		{text.FgWhite},                   // extra data
	}

	colorsWarn = colors{
		{text.BgHiYellow, text.FgHiWhite}, // prefix
		{},                                // message
		{text.FgWhite},                    // extra data
	}

	colorsError = colors{
		{text.BgHiRed, text.FgHiWhite}, // prefix
		{},                             // message
		{text.FgWhite},                 // extra data
	}

	colorsFatal = colors{
		{text.BgHiRed, text.FgBlack, text.Bold}, // prefix
		{text.BgBlack, text.FgHiRed},            // message
		{text.FgWhite},                          // extra data
	}
)
