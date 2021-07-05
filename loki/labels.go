package loki

import (
	"bufio"
	"fmt"
	"strings"
	"unicode"
)

func NewLabelParser(labels string) *LabelParser {
	return &LabelParser{
		buffer: bufio.NewReader(strings.NewReader(labels)),
		labels: make(map[string]string),
	}
}

type LabelParser struct {
	buffer *bufio.Reader
	labels map[string]string
}

func (p *LabelParser) Parse() (map[string]string, error) {
	if err := p.expect('{'); err != nil {
		return p.labels, err
	}

	if err := p.parseLabels(); err != nil {
		return p.labels, err
	}

	if err := p.expect('}'); err != nil {
		return p.labels, err
	}

	return p.labels, nil
}

func (p *LabelParser) parseLabels() error {
	for {
		if err := p.parseLabel(); err != nil {
			return err
		}

		read, _, err := p.buffer.ReadRune()
		if err != nil {
			return err
		}
		if read != ',' {
			p.buffer.UnreadRune()
			break
		}
	}

	return nil
}

func (p *LabelParser) parseLabel() error {
	if err := p.skipWhitespace(); err != nil {
		return err
	}

	name, err := p.readLabelName()
	if err != nil {
		return err
	}

	if err := p.skipWhitespace(); err != nil {
		return err
	}

	if err := p.expect('='); err != nil {
		return err
	}

	if err := p.skipWhitespace(); err != nil {
		return err
	}

	value, err := p.readLabelValue()
	if err != nil {
		return err
	}

	if err := p.skipWhitespace(); err != nil {
		return err
	}

	p.labels[name] = value

	return nil
}

func (p *LabelParser) readLabelName() (string, error) {
	name := ""

	for {
		read, _, err := p.buffer.ReadRune()
		if err != nil {
			return "", err
		}
		if p.isLabelCharacter(read) {
			name += string(read)
		} else {
			p.buffer.UnreadRune()
			break
		}
	}

	return name, nil
}

func (p *LabelParser) readLabelValue() (string, error) {
	value := ""

	if err := p.expect('"'); err != nil {
		return "", err
	}

	for {
		read, _, err := p.buffer.ReadRune()
		if err != nil {
			return "", err
		}

		if read == '"' {
			p.buffer.UnreadRune()
			break
		} else if read == '\\' {
			escaped, _, err := p.buffer.ReadRune()
			if err != nil {
				return "", err
			}

			switch {
			case escaped == '\\':
				value += "\\"
			case escaped == '"':
				value += "\""
			case escaped == 'n':
				value += "\n"
			default:
				return "", fmt.Errorf("unexpected escaped character %c, only \\\\ or \\\" or \\n is allowed", escaped)
			}

		} else {
			value += string(read)
		}
	}

	if err := p.expect('"'); err != nil {
		return "", err
	}

	return value, nil
}

func (p *LabelParser) isLabelCharacter(read rune) bool {
	switch {
	case read >= 'A' && read <= 'Z':
		return true
	case read >= 'a' && read <= 'z':
		return true
	case read >= '0' && read <= '9':
		return true
	case read == '_':
		return true
	default:
		return false
	}
}

func (p *LabelParser) skipWhitespace() error {
	for {
		read, _, err := p.buffer.ReadRune()
		if err != nil {
			return err
		}
		if !unicode.IsSpace(read) {
			p.buffer.UnreadRune()
			return nil
		}
	}
}

func (p *LabelParser) expect(symbol rune) error {
	read, _, err := p.buffer.ReadRune()
	if err != nil {
		return err
	}
	if read != symbol {
		p.buffer.UnreadRune()
		return fmt.Errorf("unexpected character %c, expected %c", read, symbol)
	}
	return nil
}
