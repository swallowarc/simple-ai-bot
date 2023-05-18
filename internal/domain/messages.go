package domain

import "fmt"

const (
	// for User
	messageFormatHelp = `【AIの使い方】
  ? : このヘルプを表示
  ?c: AI会話履歴をクリア
  ?(質問文) : AIに質問

過去%d件までの会話履歴を考慮して回答します。
話題を変える場合は会話履歴をクリアしてください。
(最後の会話から%d時間を過ぎると自動クリアされます)`
	messageFormatApproved = `AIの利用申請が承認されました。

%s`
	MessageClearHistory   = "AI会話履歴をクリアしました。"
	MessageError          = "リクエストの処理中にエラーが発生しました。しばらく待ってから再度お試しください。"
	MessageLicensePending = "AIの利用申請中です。しばらくお待ちください。"

	// for Admin
	messageFormatIssueLicense = `AIの利用申請を受け付けました。
(%s) %s`
	messageFormatSuccessApprove = "利用申請の承認が完了しました。(%s)"
	messageFormatSuccessReject  = "利用申請の棄却が完了しました。(%s)"
)

func MessageHelp() string {
	return fmt.Sprintf(messageFormatHelp, ChatHistoryLimit, ChatHistoryTTLHour)
}

func MessageIssueLicense(name, eventSourceType string) string {
	return fmt.Sprintf(messageFormatIssueLicense, eventSourceType, name)
}

func MessageApproved() string {
	return fmt.Sprintf(messageFormatApproved, MessageHelp())
}

func MessageSuccessApprove(uniqueKey string) string {
	return fmt.Sprintf(messageFormatSuccessApprove, uniqueKey)
}

func MessageSuccessReject(uniqueKey string) string {
	return fmt.Sprintf(messageFormatSuccessReject, uniqueKey)
}
