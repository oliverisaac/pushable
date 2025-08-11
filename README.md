# Pushable

This is a golang webserver that provides a simple api to trigger web pushes.

This should be implemented BEHIND A FIREWALL. This is designed for easily sending push notifications using curl.

For example:

```bash
curl -X POST -F 'topic=test-push' -F 'title=Test Push' -F 'body=This is the body of the notification' -F 'icon=https://example.com/image.png' http://push.oisaac.dev/push

```

To sign up for push notifications, you can visit push.oisaac.dev and there is a button to subscribe that device to notifications. If you visit the site from a subscribed client, you can click the "unsubscribe" button to unsubscribe.

If any field is not defined then it will default to an empty string.

# Technologies


- htmx
- golang
- javascript
