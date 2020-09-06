package job_queue

type JobQueue interface {
	CreateQueue(name string) error
	Publish(message string) error
	Consume(qName string, consumer func(msg string)) error
}

const KEY_GENERATION_EXCHANGE = "key-generation-event"
const KEY_GENERATION_QUEUE = "key-generation-queue"
const KG_PM = "GENERATE_KEY"