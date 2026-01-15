import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Rating {
  user_id: string;
  song_id: string;
  rating: number;
  created: string;
}

export interface SongRatingsResponse {
  ratings: Rating[];
  average: number;
  count: number;
}

@Injectable({ providedIn: 'root' })
export class RatingsService {
  private readonly apiBase = '/api/v1';

  constructor(private http: HttpClient) { }

  private getAuthHeaders(): HttpHeaders {
    const token = localStorage.getItem('token');
    return token ? new HttpHeaders({ Authorization: `Bearer ${token}` }) : new HttpHeaders();
  }

  // Kreiraj ili ažuriraj ocenu
  rateSong(songId: string, rating: number): Observable<Rating> {
    return this.http.post<Rating>(
      `${this.apiBase}/ratings`,
      { song_id: songId, rating },
      { headers: this.getAuthHeaders() }
    );
  }

  // Dobij sve ocene trenutnog korisnika
  getMyRatings(): Observable<Rating[]> {
    return this.http.get<Rating[]>(`${this.apiBase}/ratings`, { headers: this.getAuthHeaders() });
  }

  // Dobij sve ocene za pesmu (prosek, broj ocena)
  getSongRatings(songId: string): Observable<SongRatingsResponse> {
    return this.http.get<SongRatingsResponse>(`${this.apiBase}/ratings/${songId}`, { headers: this.getAuthHeaders() });
  }

  // Obriši ocenu
  deleteRating(songId: string): Observable<any> {
    return this.http.delete(`${this.apiBase}/ratings/${songId}`, { headers: this.getAuthHeaders() });
  }
}
