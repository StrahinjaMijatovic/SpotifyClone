import { Component, OnInit, OnDestroy, ChangeDetectorRef } from '@angular/core';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Subject, Subscription } from 'rxjs';
import { debounceTime, distinctUntilChanged } from 'rxjs/operators';

import { ContentService } from '../../services/content.service';
import { AuthService } from '../../services/auth.service';
import type { Artist, Album, Song } from '../../models/content.models';

@Component({
  selector: 'app-home-component',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './home-component.html',
  styleUrls: ['./home-component.css'],
})
export class HomeComponent implements OnInit, OnDestroy {
  searchQuery = '';
  artists: Artist[] = [];
  albums: Album[] = [];
  songs: Song[] = [];
  genres: any[] = [];
  selectedGenreId = '';
  isSearchActive = false;
  loading = false;
  errorMessage: string | null = null;

  private searchSubject = new Subject<string>();
  private searchSubscription?: Subscription;

  get isAdmin(): boolean {
    return this.authService.isAdmin();
  }

  constructor(
    private contentService: ContentService,
    private authService: AuthService,
    private router: Router,
    private cdr: ChangeDetectorRef
  ) { }

  ngOnInit(): void {
    this.loadGenres();
    this.loadArtists();
    this.setupLiveSearch();
  }

  ngOnDestroy(): void {
    this.searchSubscription?.unsubscribe();
  }

  private setupLiveSearch(): void {
    this.searchSubscription = this.searchSubject.pipe(
      debounceTime(300),
      distinctUntilChanged()
    ).subscribe(query => {
      if (query.trim()) {
        this.performSearch(query);
      } else {
        this.loadArtists();
      }
    });
  }

  onSearchInput(): void {
    this.searchSubject.next(this.searchQuery);
  }

  goToProfile(): void {
    this.router.navigate(['/profile']);
  }

  goToAdmin(): void {
    this.router.navigate(['/admin']);
  }

  goToNotifications(): void {
    this.router.navigate(['/notifications']);
  }

  goToArtist(artistId: string): void {
    this.router.navigate(['/artist', artistId]);
  }

  loadGenres(): void {
    this.contentService.getGenres().subscribe({
      next: (data) => {
        this.genres = data ?? [];
        this.cdr.detectChanges();
      },
      error: () => {
        console.error('Failed to load genres');
        this.cdr.detectChanges();
      },
    });
  }

  loadArtists(): void {
    this.loading = true;
    this.errorMessage = null;
    this.isSearchActive = false;
    this.albums = [];
    this.songs = [];

    this.contentService.getArtists(this.selectedGenreId || undefined).subscribe({
      next: (data) => {
        this.artists = data ?? [];
        this.loading = false;
        this.cdr.detectChanges();
      },
      error: () => {
        this.errorMessage = 'Ne mogu da učitam umetnike.';
        this.loading = false;
        this.cdr.detectChanges();
      },
    });
  }

  onGenreChange(): void {
    this.searchQuery = '';
    this.loadArtists();
  }

  search(): void {
    const q = this.searchQuery.trim();
    if (!q) return;
    this.performSearch(q);
  }

  private performSearch(query: string): void {
    this.loading = true;
    this.errorMessage = null;
    this.selectedGenreId = '';
    this.isSearchActive = true;

    this.contentService.search(query).subscribe({
      next: (data) => {
        this.artists = data?.artists ?? [];
        this.albums = data?.albums ?? [];
        this.songs = data?.songs ?? [];
        this.loading = false;
        this.cdr.detectChanges();
      },
      error: () => {
        this.errorMessage = 'Pretraga nije uspela.';
        this.loading = false;
        this.cdr.detectChanges();
      },
    });
  }

  clearSearch(): void {
    this.searchQuery = '';
    this.selectedGenreId = '';
    this.loadArtists();
  }

  goToAlbum(albumId: string): void {
    this.router.navigate(['/album', albumId]);
  }

  formatDuration(seconds: number): string {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  }

  logout(): void {
    localStorage.removeItem('token');
    this.router.navigate(['/login']);
  }
}
