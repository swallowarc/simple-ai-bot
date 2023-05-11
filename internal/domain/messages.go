package domain

import "fmt"

const (
	messageFormatHelp = `【AIの使い方】
  ? : このヘルプを表示
  ?c: AI会話履歴をクリア
  ?(質問文) : AIに質問

過去%d件までの会話履歴を考慮して回答します。
話題を変える場合は会話履歴をクリアしてください。
(最後の会話から%d時間を過ぎると自動クリアされます)`
	messageFormatIssueLicense = `AIの利用申請を受け付けました。
%s
%s`
	messageFormatApproved = `AIの利用申請が承認されました。

%s`

	MessageClearHistory   = "AI会話履歴をクリアしました。"
	MessageError          = "リクエストの処理中にエラーが発生しました。しばらく待ってから再度お試しください。"
	MessageLicensePending = "AIの利用申請中です。しばらくお待ちください。"
)

func MessageHelp() string {
	return fmt.Sprintf(messageFormatHelp, ChatHistoryLimit, ChatHistoryTTLHour)
}

func MessageIssueLicense(name, eventSourceID string) string {
	return fmt.Sprintf(messageFormatIssueLicense, name, eventSourceID)
}

func MessageApproved() string {
	return fmt.Sprintf(messageFormatApproved, MessageHelp())
}
