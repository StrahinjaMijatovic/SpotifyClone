import { Component, OnInit, ChangeDetectorRef } from '@angular/core';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';

import { NotificationsService, Notification } from '../../services/notifications.service';

@Component({
  selector: 'app-notifications',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './notifications.html',
  styleUrls: ['./notifications.css'],
})
export class NotificationsComponent implements OnInit {
  notifications: Notification[] = [];
  loading = true;
  errorMessage: string | null = null;

  constructor(
    private notificationsService: NotificationsService,
    private router: Router,
    private cdr: ChangeDetectorRef
  ) { }

  ngOnInit(): void {
    this.loadNotifications();
  }

  loadNotifications(): void {
    this.loading = true;
    this.errorMessage = null;

    this.notificationsService.getNotifications().subscribe({
      next: (data) => {
        this.notifications = data ?? [];
        this.loading = false;
        this.cdr.detectChanges();
      },
      error: (err) => {
        this.errorMessage = 'Ne mogu da učitam notifikacije.';
        this.loading = false;
        this.cdr.detectChanges();
      },
    });
  }

  markAsRead(notification: Notification): void {
    if (notification.read) return;

    this.notificationsService.markAsRead(notification.id).subscribe({
      next: () => {
        notification.read = true;
        this.cdr.detectChanges();
      },
      error: () => { }
    });
  }

  markAllAsRead(): void {
    const unread = this.notifications.filter(n => !n.read);
    unread.forEach(n => this.markAsRead(n));
  }

  get unreadCount(): number {
    return this.notifications.filter(n => !n.read).length;
  }

  getNotificationIcon(type: string): string {
    switch (type) {
      case 'new_album': return '💿';
      case 'new_song': return '🎵';
      case 'rating': return '⭐';
      case 'system': return '🔔';
      default: return '📩';
    }
  }

  formatDate(dateString: string): string {
    const date = new Date(dateString);
    return date.toLocaleDateString('sr-RS', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }

  goBack(): void {
    this.router.navigate(['/home']);
  }
}
