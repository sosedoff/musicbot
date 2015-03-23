package bot

import (
	"regexp"
)

type HandlerFunc func(*Bot, *Match)

type Command struct {
	re      *regexp.Regexp
	bot     *Bot
	handler HandlerFunc
}

type Match struct {
	Text   string
	Values []string
}

func NewCommand(bot *Bot, expr string, handler HandlerFunc) Command {
	return Command{
		bot:     bot,
		re:      regexp.MustCompile(expr),
		handler: handler,
	}
}

func (cmd *Command) Match(text string) *Match {
	results := cmd.re.FindStringSubmatch(text)
	if len(results) == 0 {
		return nil
	}

	return &Match{
		Text:   text,
		Values: results[1:],
	}
}
