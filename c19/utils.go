package c19

import (
	"fmt"
	"strings"

	"github.com/NicoNex/covidtron-19000/apiutil"
)

func formatNote(nota string, ntype NoteType) string {
	var msg strings.Builder

	msg.WriteString("\n\n*Note")

	switch ntype {
	case Note:
		msg.WriteString(" generali")
	case NoteCasi:
		msg.WriteString(" relative ai test effettuati")
	case NoteTest:
		msg.WriteString(" relative ai casi testati")
	}

	msg.WriteString(":*")

	note := strings.Split(nota, ". ")

	for i, n := range note {
		n = strings.TrimSuffix(n, "  ")

		if !strings.HasSuffix(n, ".") {
			n += "."
		}

		if strings.Contains(n, "  -") {
			spl := strings.Split(n, "  -")

			for _, s := range spl {
				if strings.HasPrefix(s, " ") {
					s = "-" + s
				}

				msg.WriteString(fmt.Sprintf("\n%s", s))
			}
		} else if strings.TrimSpace(n) != "." {
			if i == 0 || (i > 0 && len(note[i-1]) != 6) {
				msg.WriteString(fmt.Sprintf("\n- %s", n))
			} else {
				msg.WriteString(fmt.Sprintf(" %s", n))
			}
		}
	}

	return msg.String()
}

func plus(value int) string {
	if value > 0 {
		return "+" + apiutil.Ifmt(value)
	}
	return apiutil.Ifmt(value)
}
