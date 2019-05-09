package endpoint_poller



func CreateEndPointPoll() EndPointPoll {
     return &modBusRtuPoller{}
}