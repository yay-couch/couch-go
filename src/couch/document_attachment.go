package couch

type DocumentAttachment struct {
    Document    *Document
    file        string
    fileNam     string
    data        string
    dataLength  uint
    contentType string
    digest      string
}

func NewDocumentAttachment() *DocumentAttachment {
    return &DocumentAttachment{
        //
    }
}
