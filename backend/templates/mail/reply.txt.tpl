Hello, {{ .Email }}!

Your recipient student {{ .RecipientID }} replied to your message. To view the message click the link below:
{{ .FrontendURL }}/messages/{{ .RecipientID }}/{{ .MessageID }}

- Mr. Kupido