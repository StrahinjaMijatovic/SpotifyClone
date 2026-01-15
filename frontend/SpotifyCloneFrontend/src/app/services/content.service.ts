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

  // Admin: Genres
  getGenres(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiBase}/genres`);
  }

  // Get single artist
  getArtist(id: string): Observable<Artist> {
    return this.http.get<Artist>(`${this.apiBase}/artists/${id}`, { headers: this.getAuthHeaders() });
  }

  createArtist(data: { name: string; biography: string; genres: string[] }): Observable<Artist> {
    return this.http.post<Artist>(`${this.apiBase}/artists`, data);
  }

  updateArtist(id: string, data: { name?: string; biography?: string; genres?: string[] }): Observable<any> {
    return this.http.put(`${this.apiBase}/artists/${id}`, data);
  }

  // Get single album
  getAlbum(id: string): Observable<Album> {
    return this.http.get<Album>(`${this.apiBase}/albums/${id}`, { headers: this.getAuthHeaders() });
  }

  // Get albums by artist
  getAlbumsByArtist(artistId: string): Observable<Album[]> {
    const params = new HttpParams().set('artist_id', artistId);
    return this.http.get<Album[]>(`${this.apiBase}/albums`, { headers: this.getAuthHeaders(), params });
  }

  // Get songs by album
  getSongsByAlbum(albumId: string): Observable<Song[]> {
    const params = new HttpParams().set('album_id', albumId);
    return this.http.get<Song[]>(`${this.apiBase}/songs`, { headers: this.getAuthHeaders(), params });
  }

  // Admin: Albums
  createAlbum(data: { name: string; date: string; genre: string; artists: string[] }): Observable<Album> {
    return this.http.post<Album>(`${this.apiBase}/albums`, data);
  }

  // Admin: Songs
  createSong(data: { name: string; duration: number; genre: string; album: string; artists: string[] }): Observable<Song> {
    return this.http.post<Song>(`${this.apiBase}/songs`, data);
  }
}
