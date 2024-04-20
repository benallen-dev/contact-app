package flash

var (
	flashMessages []string
)

func init() {
	flashMessages = []string{}
}

func Queue(message ...string) {
	flashMessages = append(flashMessages, message...)
}

func Fetch() []string {
	messages := flashMessages
	flashMessages = []string{}

	return messages
}

