package eventpub

func (evp *EventPub) typeRepositoryAdded(event *RepositoryAddedEvent) error {
	return evp.sendEvent(typeRepositoryAdded, event, evp.config.GooglePubSubTopicRepository)
}
