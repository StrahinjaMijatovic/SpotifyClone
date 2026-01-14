import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import type { Album, Artist, Song, SearchResult } from '../models/content.models';

@Injectable({ providedIn: 'root' })
export class ContentService {
  private readonly apiBase = '/api/v1';

  constructor(private http: HttpClient) { }

  private getAuthHeaders(): HttpHeaders {
    const token = localStorage.getItem('token');
    return token ? new HttpHeaders({ Authorization: `Bearer ${token}` }) : new HttpHeaders();
  }

  getArtists(): Observable<Artist[]> {
    return this.http.get<Artist[]>(`${this.apiBase}/artists`, { headers: this.getAuthHeaders() });
  }

  getAlbums(): Observable<Album[]> {
    return this.http.get<Album[]>(`${this.apiBase}/albums`, { headers: this.getAuthHeaders() });
  }

  getSongs(): Observable<Song[]> {
    return this.http.get<Song[]>(`${this.apiBase}/songs`, { headers: this.getAuthHeaders() });
  }

  search(q: string): Observable<SearchResult> {
    const params = new HttpParams().set('q', q);
    return this.http.get<SearchResult>(`${this.apiBase}/search`, {
      headers: this.getAuthHeaders(),
      params,
    });
  }
}
