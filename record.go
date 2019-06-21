package ksmglog

// Record collects all fields from ksmg json
type Record struct {
	ID          int    `json:"id"`
	Time        int    `json:"time"`
	Type        string `json:"type"`
	Result      string `json:"result"`
	Person      string `json:"person"`
	Description string `json:"description"`
	EventName   string `json:"eventName"`
	Details     struct {
		MessageInfo struct {
			MessageId      string   `json:"messageId"`
			Size           string   `json:"size"`
			SmtpMessageId  string   `json:"smtpMessageId"`
			ClientAddress  string   `json:"clientAddress"`
			ClientHostName string   `json:"clientHostName"`
			From           string   `json:"from"`
			To             []string `json:"to"`
			Cc             []string `json:"cc"`
			Bcc            []string `json:"bcc"`
			Subject        string   `json:"subject"`
		} `json:"messageInfo"`
		Rules                []int  `json:"rules"`
		AvStatus             string `json:"avStatus"`
		DocWithMacroDetected bool   `json:"docWithMacroDetected"`
		AvNotScannedReason   string `json:"avNotScannedReason"`
		AsStatus             string `json:"asStatus"`
		AsNotScannedReason   string `json:"asNotScannedReason"`
		MaStatus             string `json:"maStatus"`
		MaNotScannedReason   string `json:"maNotScannedReason"`
		ApStatus             string `json:"apStatus"`
		ApNotScannedReason   string `json:"apNotScannedReason"`
		CfStatus             string `json:"cfStatus"`
		CfNotScannedReason   string `json:"cfNotScannedReason"`
		KtStatus             string `json:"ktStatus"`
		KtNotScannedReason   string `json:"ktNotScannedReason"`
		KtSkipReason         string `json:"ktSkipReason"`
		AvMessageSizeLimit   string `json:"avMessageSizeLimit"`
		AsMessageSizeLimit   string `json:"asMessageSizeLimit"`
		AsMethod             string `json:"asMethod"`
		KtProceededBy        string `json:"ktProceededBy"`
		ApuMethod            string `json:"apuMethod"`
		WmufMethod           string `json:"wmufMethod"`
		PartResults          []struct {
			FileName string `json:"fileName"`
			FileSize string `json:"fileSize"`
			AvInfo   struct {
				Statuses []struct {
					AvStatus string `json:"avStatus"`
				} `json:"statuses"`
				DocWithMacroDetected bool     `json:"docWithMacroDetected"`
				SkipReason           string   `json:"skipReason"`
				SkipDescription      string   `json:"skipDescription"`
				Threats              []string `json:"threats"`
				DisinfectedObjects   []string `json:"disinfectedObjects"`
				DeletedObjects       []string `json:"deletedObjects"`
			} `json:"avInfo"`
			CfInfo struct {
				Statuses         []string `json:"statuses"`
				BannedFileName   string   `json:"bannedFileName"`
				BannedFileFormat string   `json:"bannedFileFormat"`
			} `json:"cfInfo"`
			Action string `json:"action"`
		} `json:"partResults"`
		MaInfo struct {
			DmarcVerdict string   `json:"dmarcVerdict"`
			SpfVerdict   string   `json:"spfVerdict"`
			DkimVerdicts []string `json:"dkimVerdicts"`
		} `json:"maInfo"`
		Action                       string   `json:"action"`
		BackupReason                 string   `json:"backupReason"`
		UnsafeNotificationRecipients []string `json:"unsafeNotificationRecipients"`
	} `json:"details"`
}
