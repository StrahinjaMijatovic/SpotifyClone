import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Notification {
  id: string;
  user_id: string;
  message: string;
  type: string;
  read: boolean;
  created_at: string;
}

@Injectable({ providedIn: 'root' })
export class NotificationsService {
  private readonly apiBase = '/api/v1';

  constructor(private http: HttpClient) { }

  private getAuthHeaders(): HttpHeaders {
    const token = localStorage.getItem('token');
    return token ? new HttpHeaders({ Authorization: `Bearer ${token}` }) : new HttpHeaders();
  }

  // Dobij sve notifikacije korisnika
  getNotifications(): Observable<Notification[]> {
    return this.http.get<Notification[]>(`${this.apiBase}/notifications`, { headers: this.getAuthHeaders() });
  }

  // Označi notifikaciju kao pročitanu
  markAsRead(id: string): Observable<any> {
    return this.http.put(`${this.apiBase}/notifications/${id}/read`, {}, { headers: this.getAuthHeaders() });
  }

  getUnreadCount(): Observable<{ unread_count: number }> {
    return this.http.get<{ unread_count: number }>(`${this.apiBase}/notifications/unread/count`, { headers: this.getAuthHeaders() });
  }
}
