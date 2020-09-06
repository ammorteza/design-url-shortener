package job_queue

type JobQueue interface {
	CreateQueue(name string) error
	Publish(id string, body []byte) error
	Consume(qName string, consumer func(id string, body []byte)) error
}

const URL_SHORTENER_EXCHANGE = "url-shortener-event"
const URL_SHORTENER_QUEUE = "url-shortener-queue"
const INCREASE_VISIT_COUNT_PM = "increase-visit-count"