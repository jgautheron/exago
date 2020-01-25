package eventpub

func (evp *EventPub) typeRepositoryAdded(event *RepositoryAddedEvent) error {
	return evp.sendEvent(TypeRepositoryAdded, event, evp.config.GooglePubSubTopicRepository)
}
