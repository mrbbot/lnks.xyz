export const notificationMixin = {
  data: function() {
    return {
      notifications: [],
      totalNotificationCount: 0
    };
  },
  methods: {
    showNotification: function(content) {
      const id = this.totalNotificationCount;
      this.totalNotificationCount++;
      this.notifications.push({id, content});
      setTimeout(() => {
        this.notifications = this.notifications.filter(n => n.id !== id);
      }, 3000);
    }
  },
  mounted: function() {
    for(const notification of window.notifcations) {
      this.showNotification(notification);
    }
  }
};
