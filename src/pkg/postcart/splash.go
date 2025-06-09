package postcart

import (
	"fmt"

	"github.com/agladfield/postcart/pkg/shared/tools/splash"
)

const postcartSplashHeader = `                 _                        _
                | |                      | |
 _ __   ___  ___| |_ ___       __ _ _ ___| |__,
| '_ \ / _ \/ __| __/ __|     / _  | '__/|  __/
| |_) | (_) \__ \ || (_(     |_| | | |   | |
| .__/ \___/|___/\__\___| /\  \__,_|_|    \__\
| |                       \/
|_|                                     `

const postcartSplashCenter = ` A Postmark Inbox Innovators Challenge Entry
         dev.to/challenges/postmark
#devchallenge #postmarkchallenge #webdev #ai`

const postcartSplashFooter = `By Arthur Gladfield
    @agladfield`

func postcartSplash() {
	str := splash.Splash(splash.SplashContentOptions{
		Header: postcartSplashHeader,
		Center: postcartSplashCenter,
		Footer: postcartSplashFooter,
	}, splash.SplashConfigOptions{
		BorderChar: "@",
	})
	fmt.Println(str)
}
