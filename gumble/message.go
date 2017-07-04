package gumble // import "github.com/101Bas/gumble/gumble"

// Message is data that be encoded and sent to the server. The following
// types implement this interface:
//  AccessTokens
//  ACL
//  BanList
//  RegisteredUsers
//  TextMessage
//  VoiceTarget
type Message interface {
	writeMessage(client *Client) error
}
