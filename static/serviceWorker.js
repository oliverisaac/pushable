self.addEventListener('install', function(event) {
  self.skipWaiting();
});

self.addEventListener('activate', function(event) {
  event.waitUntil(clients.claim());
});

self.addEventListener('push', function(event) {
    const data = event.data.json();
    event.waitUntil(
        self.registration.showNotification(data.title, {
            body: data.body,
            icon: data.icon,
            badge: data.badge,
            data: data.data
        })
    );
});

self.addEventListener('notificationclick', function(event) {
  event.notification.close();

  event.waitUntil(
    clients.matchAll({ type: 'window', includeUncontrolled: true }).then(function(clientList) {
      if (event.notification.data.link === undefined ) {
        return
      }

      const paramsData = {
        target: event.notification.data.link,
      };

      const params = new URLSearchParams(paramsData);
      const redirectURL = "/redirect?" + params.toString();
      if (clientList.length > 0) {
        let client = clientList[0];
        for (let i = 0; i < clientList.length; i++) {
          if (clientList[i].focused) {
            client = clientList[i];
          }
        }
        client.navigate(redirectURL);
        return client.focus();
      }
      return clients.openWindow(redirectURL);
    })
  );
});