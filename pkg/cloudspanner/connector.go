package cloudspanner

type CloudSpannerConnector struct {
}

func NewCloudSpannerConnector() *CloudSpannerConnector {
	return &CloudSpannerConnector{}
}

func (csc *CloudSpannerConnector) Read()  {}
func (csc *CloudSpannerConnector) Write() {}
