package eventpub

func (evp *EventPub) RepositoryAdded(event *RepositoryAddedEvent) error {
	return evp.sendEvent(TypeRepositoryAdded, event, evp.config.GooglePubSubTopicRepository)
}
